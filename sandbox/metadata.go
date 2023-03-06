/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package sandbox

import (
	"context"

	"github.com/containerd/ttrpc"
	"google.golang.org/grpc/metadata"
)

// NetNSPathKey is the key of sandbox network namespace path
const (
	NetNSPathKey = "io.containerd.sandbox.netns_path"
)

// WithCustomKeyValue sets a custom key value pair on the context
func WithCustomKeyValue(ctx context.Context, key, value string) context.Context {
	return withTTRPCKeyValueHeader(withGRPCKeyValueHeader(ctx, key, value), key, value)
}

func withGRPCKeyValueHeader(ctx context.Context, key, value string) context.Context {
	// also store on the grpc headers so it gets picked up by any clients that
	// are using this.
	nsheader := metadata.Pairs(key, value)
	md, ok := metadata.FromOutgoingContext(ctx) // merge with outgoing context.
	if !ok {
		md = nsheader
	} else {
		// order ensures the latest is first in this list.
		md = metadata.Join(nsheader, md)
	}

	return metadata.NewOutgoingContext(ctx, md)
}

func withTTRPCKeyValueHeader(ctx context.Context, key, value string) context.Context {
	md, ok := ttrpc.GetMetadata(ctx)
	if !ok {
		md = ttrpc.MD{}
	} else {
		md = copyMetadata(md)
	}
	md.Set(key, value)
	return ttrpc.WithMetadata(ctx, md)
}

func copyMetadata(src ttrpc.MD) ttrpc.MD {
	md := ttrpc.MD{}
	for k, v := range src {
		md[k] = append(md[k], v...)
	}
	return md
}

func MetadataInCtx(ctx context.Context, key string) (string, bool) {
	value, ok := fromGRPCHeaderByKey(ctx, key)
	if !ok {
		return fromTTRPCHeaderByKey(ctx, key)
	}
	return value, ok
}

func fromGRPCHeaderByKey(ctx context.Context, key string) (string, bool) {
	// try to extract for use in grpc servers.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	values := md[key]
	if len(values) == 0 {
		return "", false
	}

	return values[0], true
}

func fromTTRPCHeaderByKey(ctx context.Context, key string) (string, bool) {
	return ttrpc.GetMetadataValue(ctx, key)
}
