// Copyright 2017 GRAIL, Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package values

import (
	"github.com/grailbio/base/digest"
	"github.com/grailbio/reflow/types"
)

// Symtab is a symbol table of values.
type Symtab map[string]T

// Env binds identifiers to evaluation.
type Env struct {
	// Symtab is the symbol table for this level.
	Symtab Symtab
	debug  bool
	next   *Env
}

// NewEnv constructs and initializes a new Env.
func NewEnv() *Env {
	var e *Env
	return e.Push()
}

// Bind binds the identifier id to value v.
func (e *Env) Bind(id string, v T) {
	e.Symtab[id] = v
}

type digester interface {
	Digest() digest.Digest
}

// Digest returns the digest for the value with idenfier id.
// The supplied type is used to compute the digest. If the value
// is a digest.Digest, it is returned directly; if implements the
// interface
//
//	interface{
//		Digest() digest.Digest
//	}
//
// it returns the result of calling the Digest Method.
func (e *Env) Digest(id string, t *types.T) digest.Digest {
	v := e.Value(id)
	switch vv := v.(type) {
	case digest.Digest:
		return vv
	case digester:
		return vv.Digest()
	default:
		return Digest(vv, t)
	}
}

// Level returns the level of identifier id. Level can thus be used
// as a de-Bruijn index (in conjunction with the identifier).
func (e *Env) Level(id string) int {
	for l := 0; e != nil; e, l = e.next, l+1 {
		if e.Symtab[id] != nil {
			return l
		}
	}
	return -1
}

// Contains tells whether environment e binds identifier id.
func (e *Env) Contains(id string) bool {
	for ; e != nil; e = e.next {
		if _, ok := e.Symtab[id]; ok {
			return true
		}
	}
	return false
}

// Value returns the value bound to identifier id, or else nil.
func (e *Env) Value(id string) T {
	if e == nil {
		return nil
	}
	for ; e != nil; e = e.next {
		if v := e.Symtab[id]; v != nil {
			return v
		}
	}
	return nil
}

// Push returns returns a new environment level, linked
// to the previous.
func (e *Env) Push() *Env {
	return &Env{
		Symtab: make(Symtab),
		next:   e,
	}
}

// Debug sets the debug flag on this environment.
// This causes IsDebug to return true.
func (e *Env) Debug() {
	e.debug = true
}

// IsDebug returns true if the debug flag is set in this environment.
func (e *Env) IsDebug() bool {
	if e == nil {
		return false
	}
	for ; e != nil; e = e.next {
		if e.debug {
			return true
		}
	}
	return false
}
