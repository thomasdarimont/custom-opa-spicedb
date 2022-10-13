package builtins

import (
	"errors"
	"fmt"
	authzedpb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	"github.com/thomasdarimont/custom-opa/custom-opa-spicedb/plugins/authzed"
	"strings"
)

var checkPermissionBuiltinDecl = &rego.Function{
	Name: "authzed.check_permission",
	Decl: types.NewFunction(
		types.Args(types.S, types.S, types.S), // subject, permission, resource
		types.B),                              // Returns a boolean
}

// Use a custom cache key type to avoid collisions with other builtins caching data!!
type checkPermissionCacheKeyType string

// checkPermissionBuiltinImpl checks the given permission requests against spicedb.
func checkPermissionBuiltinImpl(bctx rego.BuiltinContext, subjectTerm, permissionTerm, resourceIdTerm *ast.Term) (*ast.Term, error) {

	// repository:authzed_go
	var resource string
	if err := ast.As(resourceIdTerm.Value, &resource); err != nil {
		return nil, err
	}

	// clone
	var permission string
	if err := ast.As(permissionTerm.Value, &permission); err != nil {
		return nil, err
	}

	// user:jake#...
	var subject string
	if err := ast.As(subjectTerm.Value, &subject); err != nil {
		return nil, err
	}

	// Check if it is already cached, assume they never become invalid.
	var cacheKey = checkPermissionCacheKeyType(fmt.Sprintf("%s#%s@%s", subject, permission, resource))
	cached, ok := bctx.Cache.Get(cacheKey)
	if ok {
		return ast.NewTerm(cached.(ast.Value)), nil
	}

	objectType, objectId, subjectFound := strings.Cut(subject, ":")
	if !subjectFound {
		return nil, errors.New("could not parse authzdb subject")
	}

	subjectReference := &authzedpb.SubjectReference{Object: &authzedpb.ObjectReference{
		ObjectType: objectType,
		ObjectId:   objectId,
	}}

	resourceType, resourceId, resourceFound := strings.Cut(resource, ":")
	resourceReference := &authzedpb.ObjectReference{
		ObjectType: resourceType,
		ObjectId:   resourceId,
	}

	if !resourceFound {
		return nil, errors.New("could not parse authzdb resource")
	}

	client := authzed.GetAuthzedClient()
	if client == nil {
		return nil, errors.New("authzed client not configured")
	}

	resp, err := client.CheckPermission(bctx.Context, &authzedpb.CheckPermissionRequest{
		Resource:   resourceReference,
		Permission: permission,
		Subject:    subjectReference,
	})

	if err != nil {
		return nil, err
	}

	result := ast.Boolean(resp.Permissionship == authzedpb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION)
	bctx.Cache.Put(cacheKey, result)

	return ast.NewTerm(result), nil
}
