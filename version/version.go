// Copyright 2018 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package version provides build version information.
package version

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// The following fields are populated at build time using -ldflags -X.
// Note that DATE is omitted for reproducible builds
var (
	buildVersion     = "unknown"
	buildGitRevision = "unknown"
	buildStatus      = "unknown"
	buildTag         = "unknown"
	buildHub         = "unknown"
	buildVendor      = "oss"
)

// BuildInfo describes version information about the binary build.
type BuildInfo struct {
	Version       string `json:"version"`
	GitRevision   string `json:"revision"`
	GolangVersion string `json:"golang_version"`
	BuildStatus   string `json:"status"`
	GitTag        string `json:"tag"`
	// Vendor denotes who built the image. Default is "istio.io".
	Vendor        string `json:"vendor"`
}

// ServerInfo contains the version for a single control plane component
type ServerInfo struct {
	Component string
	Info      BuildInfo
}

// MeshInfo contains the versions for all Istio control plane components
type MeshInfo []ServerInfo

// ProxyInfo contains the version for a single data plane component
type ProxyInfo struct {
	ID           string
	IstioVersion string
}

// DockerBuildInfo contains and exposes Hub: buildHub and Tag: buildVersion
type DockerBuildInfo struct {
	Hub string
	Tag string
}

// NewBuildInfoFromOldString creates a BuildInfo struct based on the output
// of previous Istio components '-- version' output
func NewBuildInfoFromOldString(oldOutput string) (BuildInfo, error) {
	res := BuildInfo{}

	lines := strings.Split(oldOutput, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.SplitN(line, ":", 2)
		if fields != nil {
			if len(fields) != 2 {
				return BuildInfo{}, fmt.Errorf("invalid BuildInfo input, field '%s' is not valid", fields[0])
			}
			value := strings.TrimSpace(fields[1])
			switch fields[0] {
			case "Version":
				res.Version = value
			case "GitRevision":
				res.GitRevision = value
			case "GolangVersion":
				res.GolangVersion = value
			case "BuildStatus":
				res.BuildStatus = value
			case "GitTag":
				res.GitTag = value
			default:
				// Skip unknown fields, as older versions may report other fields
				continue
			}
		}
	}

	return res, nil
}

var (
	// Info exports the build version information.
	Info       BuildInfo
	DockerInfo DockerBuildInfo
)

// String produces a single-line version info
//
// This looks like:
//
// ```
// u<version>-<git revision>-<build status>
// ```
func (b BuildInfo) String() string {
	return fmt.Sprintf("%v-%v-%v",
		b.Version,
		b.GitRevision,
		b.BuildStatus)
}

// adjustCommand returns the last component of the
// OS-specific command path for use in User-Agent.
func adjustCommand(p string) string {
	// Unlikely, but better than returning "".
	if len(p) == 0 {
		return "unknown"
	}
	return filepath.Base(p)
}

// UserAgent produces a UserAgent string that can be used to self
// identify the istio binary.
//
// ```
// u<binary_name>/<version> ($GOOS/$GOARCH) <vendor>
// istioctl/1.11.2 (linux/amd64) istio/oss
// istioctl/1.11.6 (darwin/amd64) istio/Google
// ```
func (b BuildInfo) UserAgent() string {
	return fmt.Sprintf("%v/%v (%v/%v) istio/%v",
		adjustCommand(os.Args[0]),
		b.Version,
		runtime.GOOS,
		runtime.GOARCH,
		b.Vendor)
}

// LongForm returns a dump of the Info struct
// This looks like:
//
func (b BuildInfo) LongForm() string {
	return fmt.Sprintf("%#v", b)
}

func init() {
	Info = BuildInfo{
		Version:       buildVersion,
		GitRevision:   buildGitRevision,
		GolangVersion: runtime.Version(),
		BuildStatus:   buildStatus,
		GitTag:        buildTag,
		Vendor:        buildVendor,
	}

	DockerInfo = DockerBuildInfo{
		Hub: buildHub,
		Tag: buildVersion,
	}
}
