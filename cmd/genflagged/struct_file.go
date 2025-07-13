package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"
)

func (f *File) isValidStructFile() bool {
	return len(f.flagValues) != 0
}

// genStructDecl processes one 'type <name> struct' declaration clause.
// Its target fields are these whose type is bool or an alias to bool,
// and aren't embedded fields nor has the name '_'.
// Note: it doesn't include fields whose type is named bool-based types,
// as the implementation doesn't support converting to/from bool-based
// types yet.
// TODO: maybe add support for that?
func (f *File) genStructDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.TYPE {
		// We only care about 'type' declarations.
		return true
	}

	// Loop over the elements of the declaration.
	// Each element is a TypeSpec, which is a type declaration (struct, int, etc.),
	// but we are only interested in struct types.
	for _, spec := range decl.Specs {
		// Guaranteed to succeed as this is TYPE.
		tspec := spec.(*ast.TypeSpec)

		// Skip if this is not the type we're looking for.
		if f.sourceTypeName != tspec.Name.Name {
			continue
		}

		// Mark the current source type as found, regardless it's of correct
		// type or not, and regardless we will find matching fields or not.
		foundSourceType, ok := f.pkg.defs[tspec.Name]
		if !ok {
			log.Fatalf("error: no type definition found for type %s", tspec.Name)
		}
		f.foundSourceType = foundSourceType
		verbose.Printf("info: found matching type %s\n", foundSourceType)

		// Skip if this is not a struct type.
		stype, ok := tspec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		// Init the fields list, assuming the struct contains only target fields,
		// and only one field per declaration.
		f.flagValues = make([]flagValue, 0, len(stype.Fields.List))

		verbose.Printf(
			"info: proccessing %d field declarations for type %s\n",
			len(stype.Fields.List),
			tspec.Name.Name,
		)

		// Loop over the list of fields, filter only target fields.
		for idx, field := range stype.Fields.List {
			verbose.Printf(
				"info: proccessing field declaration at index %d for type %s with %d names\n",
				idx,
				tspec.Name.Name,
				len(field.Names),
			)

			// Skip embedded types.
			if len(field.Names) == 0 {
				continue
			}

			// Loop over each name in the same field declaration.
			for _, name := range field.Names {
				verbose.Printf(
					"info: proccessing field name %s for type %s\n",
					name.Name,
					tspec.Name.Name,
				)

				// Skip fields named '_', regardless of type.
				if name.Name == "_" {
					continue
				}

				// This dance lets the type checker find the aliased type
				// for us, if any.
				// It's a bit tricky: look up the object declared by name,
				// unalias its type, and use that type to do the checks.
				obj, ok := f.pkg.defs[name]
				if !ok {
					// TODO: can this ever happen??
					//  as the file, which this type is defined in, belongs
					//  to the same package we are querying this type for.
					log.Fatalf(
						"error: no field definition found for field %s from type %s",
						name,
						tspec.Name.Name,
					)
				}

				// Get the actual type of the field.
				// Note: it must be a builtin bool or an alias to one.
				actualType := types.Unalias(obj.Type())

				// Skip named types, as they're not supported.
				basicType, ok := actualType.(*types.Basic)
				if !ok {
					verbose.Printf(
						"info: found field %s but with unsupported actual type %T in type %s\n",
						name.Name,
						actualType,
						tspec.Name.Name,
					)

					continue
				}

				// Skip this field if type isn't bool.
				info := basicType.Info()
				if info&types.IsBoolean == 0 {
					verbose.Printf(
						"info: found field %s but with non-boolean actual type %T in type %s\n",
						name.Name,
						actualType,
						tspec.Name.Name,
					)

					continue
				}

				// TODO: maybe add some validation to make sure the generated types and flags
				// doesn't already exist in the package, since we have the type info about it.
				fv := flagValue{
					Field: name.Name,
					Flag:  flagName(name.Name, f.pkg.trimPrefix, f.pkg.trimSuffix),
				}
				f.flagValues = append(f.flagValues, fv)

				verbose.Printf(
					"info: added flag %s for field %s from type %s with total %d flags\n",
					fv.Flag,
					fv.Field,
					tspec.Name.Name,
					len(f.flagValues),
				)
			}
		}
	}

	// Set the flags size based on the number of loaded flag values.
	f.flagsSize = flagSize(len(f.flagValues))

	verbose.Printf(
		"info: type size is %d for type %s with total %d flags\n",
		f.flagsSize,
		f.sourceTypeName,
		len(f.flagValues),
	)

	return false
}
