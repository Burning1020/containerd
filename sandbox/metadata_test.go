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
	"testing"
)

func TestMetadataInCtx(t *testing.T) {
	ctx := context.Background()
	key := "test-key"
	value, ok := MetadataInCtx(ctx, key)
	if ok {
		t.Fatal("key should not be present in context")
	}

	if value != "" {
		t.Fatalf("value should not be defined: got %q", value)
	}

	expected := "test-value"
	nctx := WithCustomKeyValue(ctx, key, expected)

	value, ok = MetadataInCtx(nctx, key)
	if !ok {
		t.Fatal("expected to find a namespace")
	}

	if value != expected {
		t.Fatalf("unexpected namespace: %q != %q", value, expected)
	}
}
