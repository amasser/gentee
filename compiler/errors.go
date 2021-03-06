// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gentee/gentee/core"
)

const (
	// The list of errors

	// ErrSuccess means no error
	ErrSuccess = iota
	// ErrDecl is returned when the unexpexted token has been found on the top level
	ErrDecl
	// ErrLCurly is returned when the unexpexted token, expecting {
	ErrLCurly
	// ErrEnd is returned when the unexpexted end of the source, expecting }
	ErrEnd
	// ErrExp is returned when the unexpected token, expecting expression or statement {
	ErrExp
	// ErrName is return when compiler is expecting the identifier
	ErrName
	// ErrValue is returned when the unexpected token, expecting value, identifier or calling func
	ErrValue
	// ErrRun is returned when the compiler has found the second run function.
	ErrRun
	// ErrType is returned when the unexpected token, expecting type name
	ErrType
	// ErrReturn is returned when the function returns a value but it must not return
	ErrReturn
	// ErrMustReturn is returned when the function doesn't return a value but it must return
	ErrMustReturn
	// ErrReturnType is returned when the function returns a wrong type
	ErrReturnType
	// ErrOutOfRange is returned when the number is out of range
	ErrOutOfRange
	// ErrLPar is returned when there is an unclosed left parenthesis
	ErrLPar
	// ErrRPar is returned when extra right parenthesis has been found
	ErrRPar
	// ErrLSBracket is returned when there is an unclosed left square bracket
	ErrLSBracket
	// ErrRSBracket is returned when extra right square bracket has been found
	ErrRSBracket
	// ErrEmptyCode is returned when the source code is empty
	ErrEmptyCode
	// ErrFunction is returned when the compiler could not find a corresponding function
	ErrFunction
	// ErrBoolExp is returned when the compiler expects boolean result but gets different type
	ErrBoolExp
	// ErrFuncExists is returned when the function ahs already been defined
	ErrFuncExists
	// ErrUsedName is returned when the specified name has already been used
	ErrUsedName
	// ErrUnknownIdent is returned when the compiler gets unknown identifier
	ErrUnknownIdent
	// ErrLValue is returned when left operand of assign is not l-value
	ErrLValue
	// ErrOper is return when there is not operator
	ErrOper
	// ErrBoolOper is returned when && or || gets not boolen operands
	ErrBoolOper
	// ErrQuestion is returned when exp1 and exp2 have different types
	ErrQuestion
	// ErrQuestionPars is returned when ?(condition,exp1,exp2) has wrong parameters
	ErrQuestionPars
	// ErrCapitalLetters is returned when the var or func name consists of only capital letters
	ErrCapitalLetters
	// ErrConstName is returned when the name of constant doesn't consist of only capital letters
	ErrConstName
	// ErrMustAssign is returned when the constant is described without assign
	ErrMustAssign
	// ErrIota is returned when IOTA is used outside const expression
	ErrIota
	// ErrIntOper is returned when ++ or -- gets not int value
	ErrIntOper
	// ErrDoubleQuotes is returned when there is a wrong command of backslash in double quotes strings
	ErrDoubleQuotes
	// ErrLink is returned when the script with the same name but different path is already linked
	//	ErrLink
	// ErrConstDef is returned when the constant is redefined
	ErrConstDef
	// ErrChar means that the char literal has wrong format
	ErrChar
	// ErrNoIndex means that there is not a value for index
	ErrNoIndex
	// ErrVarIndex means that there is not a variable for indexing
	ErrVarIndex
	// ErrSupportIndex means that the type of the variable doesn't support indexing
	ErrSupportIndex
	// ErrTypeIndex means that the type of the index value is wrong
	ErrTypeIndex
	// ErrForIn is returned when 'in' is missing in for statement
	ErrForIn
	// ErrIdent is returned when the name contains a dot
	ErrIdent
	// ErrWrongType is returned when we get wrong type
	ErrWrongType
	// ErrNotKeyValue is returned when initialization of map value without
	ErrNotKeyValue
	// ErrKeyValue means that key:value is used outside of map initialization
	ErrKeyValue
	// ErrLineRCurly is returned when the unexpected token, expecting a new line or }
	ErrLineRCurly
	// ErrStructField is returned when the field with this name has already been defined
	ErrStructField
	// ErrTypeExists means that the type has already been defined
	ErrTypeExists
	// ErrStructType is returned when getting field of no struct type
	ErrStructType
	// ErrStruct is returned when struct type doesn't have such field
	ErrStruct
	// ErrStructAssign is returned if structures have different types in assign expression
	ErrStructAssign
	// ErrInitField means that wrong token specified as a field name of struct
	ErrInitField
	// ErrWrongField is returned when unknown field name has been specified
	ErrWrongField
	// ErrBreak returns when break is placed outside of loops
	ErrBreak
	// ErrContinue returns when continue is placed outside of loops
	ErrContinue
	// ErrNotRPar is returned when the compiler gets unexpected token instead of )
	ErrNotRPar
	// ErrNotCase is returned if case missing after switch
	ErrNotCase
	// ErrSwitchType means that type of switch statement is wrong
	ErrSwitchType
	// ErrString is returned when expecting a string value
	ErrString
	// ErrIncludeFile is returned when an include file is incorrect
	ErrIncludeFile
	// ErrDupObject is returned when the duplicate object has been found in include or import
	ErrDupObject
	// ErrImportStr is returned when the string contains an expression
	ErrImportStr
	// ErrNewLine is retuned is case of unexpected token instead of a new line
	ErrNewLine
	// ErrAddrFunc means wrong definition of function address
	ErrAddrFunc
	// ErrNoFuncType is returned when the type is not a fn type
	ErrNoFuncType
	// ErrFnReturn is returned when function and fn type return different types
	ErrFnReturn
	// ErrFnCall is returned when fn var is called with wrong parameters
	ErrFnCall
	// ErrOptional means that the optional variable is defined in the wrong place
	ErrOptional
	// ErrFuncOptional is returned when it doesn't have such optional variable
	ErrFuncOptional
	// ErrFnOptional is returned when fn variable use optional variables
	ErrFnOptional
	// ErrTypeOptional is returned when the optional variable has wrong type
	ErrTypeOptional
	// ErrTwiceOptional is returned when the optional variable is defined more than one time
	ErrTwiceOptional
	// ErrEndOptional is returned when optional parameters are not at the end`
	ErrEndOptional
	// ErrLocalName is returned when such local name has already existed
	ErrLocalName
	// ErrLocalVariadic is returned when there is a variadic parameter in local function
	ErrLocalVariadic
	// ErrGoParam is returned when there is an unnamed parameter in go statement
	ErrGoParam
	// ErrCatch is returned if catch missing after try
	ErrCatch
	// ErrRecover returns when recover is placed outside of catch
	ErrRecover
	// ErrRetry returns when recover is placed outside of catch
	ErrRetry
	// ErrLinkIndex is returned when linker gets incorrect unit index
	ErrLinkIndex
	// ErrFnBuildIn is returned when fn variable assigned to build-in function
	ErrFnBuildIn
	// ErrFnVariadic is returned when fn variable assigned to variadic function
	ErrFnVariadic

	// ErrCompiler error. It means a bug.
	ErrCompiler

	// ErrLetter is returned when an unknown character has been found
	ErrLetter = 0x100
	// ErrWord is returned when a sequence of characters is wrong
	ErrWord = 0x200
	// ErrEnvName is returned when a environment name ${NAME} is wrong
	ErrEnvName = 0x300
	// ErrDoubleColon is returned where there are two colons in one line
	ErrDoubleColon = 0x500
)

var (
	errText = map[int]string{
		ErrLetter:      `unknown character`,
		ErrWord:        `wrong sequence of characters`,
		ErrEnvName:     `wrong environment name, expecting ${NAME}`,
		ErrDoubleColon: `colon has already been specified in this line`,

		ErrLCurly:         `unexpected token, expecting {`,
		ErrEnd:            `unexpected end of the source`,
		ErrDecl:           `expected declaration: func, run etc`,
		ErrExp:            `unexpected token, expecting expression or statement`,
		ErrName:           `unexpected token, expecting the name of the identifier`,
		ErrRun:            `run function has already been defined`,
		ErrValue:          `unexpected token, expecting value, identifier or calling func`,
		ErrType:           `unexpected token, expecting type`,
		ErrReturn:         `function cannot return any value`,
		ErrMustReturn:     `function must return a value`,
		ErrReturnType:     `function returns wrong type`,
		ErrOutOfRange:     `the number %s is out of range`,
		ErrLPar:           `there is an unclosed left parenthesis`,
		ErrRPar:           `extra right parenthesis`,
		ErrLSBracket:      `there is an unclosed left square bracket`,
		ErrRSBracket:      `extra right square bracket`,
		ErrEmptyCode:      `source code is empty`,
		ErrFunction:       `function %s%s has not been found`,
		ErrBoolExp:        `wrong type of expression, expecting boolean type`,
		ErrFuncExists:     `function %s%s has already been defined`,
		ErrUsedName:       `"%s" has already been used as the name of the function, type or variable`,
		ErrUnknownIdent:   `unknown identifier %s`,
		ErrLValue:         `expecting l-value in the left operand of assign operator`,
		ErrOper:           `unexpected token, expecting operator`,
		ErrBoolOper:       `wrong type of operands, expecting boolean type`,
		ErrQuestion:       `different types of exp1 and exp2 in ?(cond, exp1, exp2)`,
		ErrQuestionPars:   `operator ? must be called as ?(boolean condition, exp1, exp2)`,
		ErrCapitalLetters: `The name of variable, type or function can't consists of only capital letters`,
		ErrConstName:      `The name of constant must consist of only capital letters`,
		ErrMustAssign:     `unexpected token, expecting =`,
		ErrIota:           `IOTA can be only used in const expression`,
		ErrIntOper:        `wrong type of operands, expecting int type`,
		ErrDoubleQuotes:   `invalid syntax of double quotes string`,
		//		ErrLink:           `script %s has already been linked`,
		ErrConstDef:      `constant %s has already been defined!`,
		ErrChar:          `char literal has wrong format`,
		ErrNoIndex:       `there is not index value`,
		ErrVarIndex:      `unexpected token, expecting a variable for indexing`,
		ErrSupportIndex:  `%s type does not support indexing`,
		ErrTypeIndex:     `wrong type of index, expecting %s type`,
		ErrForIn:         `unexpected token, expecting 'in'`,
		ErrIdent:         `the name of the identifier can't contain a dot`,
		ErrWrongType:     `wrong type, expecting %s type`,
		ErrNotKeyValue:   `unexpected type, expecting a pair of key and value`,
		ErrKeyValue:      `unexpected a pair of key and value, expecting %s type`,
		ErrLineRCurly:    `unexpected token, expecting a new line or }`,
		ErrStructField:   `%s field has already been defined`,
		ErrTypeExists:    `%s type has already been defined`,
		ErrStructType:    `%s type is not struct type`,
		ErrStruct:        `%s type doesn't have %s field`,
		ErrStructAssign:  `can't assign %s to %s`,
		ErrInitField:     `unexpected token, expecting the name of the field`,
		ErrWrongField:    `there is not %s field in %s struct`,
		ErrBreak:         `break can only be inside while or for`,
		ErrContinue:      `continue can only be inside while or for`,
		ErrNotRPar:       `unexpected token, expecting ')'`,
		ErrNotCase:       `unexpected token, expecting 'case'`,
		ErrSwitchType:    `wrong type %s for switch, expecting int, float, bool, char or str`,
		ErrString:        `unexpected token, expecting a string`,
		ErrIncludeFile:   `can't read include file: %s`,
		ErrDupObject:     `duplicate of %s has been found after include/import`,
		ErrImportStr:     `string cannot contain an expression`,
		ErrNewLine:       `unexpected token, expecting a new line`,
		ErrAddrFunc:      `address of function must be defined as &name.type`,
		ErrNoFuncType:    `type %s is not a fn type`,
		ErrFnReturn:      `function %s and %s fn type return different types`,
		ErrFnCall:        `fn type %s is different from %s`,
		ErrOptional:      `optional variable cannot be inside 'run' or nested block`,
		ErrFuncOptional:  `function doesn't have optional variable %s`,
		ErrFnOptional:    `fn variables can't use optional variables`,
		ErrTypeOptional:  `%s optional variable has different type`,
		ErrTwiceOptional: `%s optional variable is defined more than one time`,
		ErrEndOptional:   `optional parameters must be at the end`,
		ErrLocalName:     `%s local function has already been defined`,
		ErrLocalVariadic: `local function cannot have a variadic parameter`,
		ErrGoParam:       `there is an unnamed parameter in go statement`,
		ErrCatch:         `unexpected token, expecting 'catch'`,
		ErrRecover:       `'recover' can only be inside catch`,
		ErrRetry:         `'retry' can only be inside catch`,
		ErrLinkIndex:     `incorrect link index %d`,
		ErrFnBuildIn:     `fn variable can't be assigned to a built-in function`,
		ErrFnVariadic:    `fn variable can't be assigned to a variadic function`,

		ErrCompiler: `you have found a compiler bug [%s]. Let us know, please`,
	}
)

func (cmpl *compiler) ErrorPos(pos int, errID int, pars ...interface{}) error {
	lex := cmpl.unit.Lexeme
	line, column := lex.LineColumn(pos)
	return errors.New(core.ErrFormat(lex.Path, line, column, fmt.Sprintf(errText[errID], pars...)))
}

func (cmpl *compiler) Error(errID int, pars ...interface{}) error {
	return cmpl.ErrorPos(cmpl.pos, errID, pars...)
}

func (cmpl *compiler) ErrorFunction(errID int, pos int, name string, pars []*core.TypeObject) error {
	var params []string
	for _, par := range pars {
		if par != nil {
			params = append(params, par.GetName())
		} else {
			params = append(params, `...`)
		}
	}
	return cmpl.ErrorPos(pos, errID, name, fmt.Sprintf(`(%s)`, strings.Join(params, `, `)))
}
