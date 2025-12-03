package aql

import (
	"errors"
	"unicode/utf8"
)

var (
	ErrUnexpectedEndOfInput = errors.New("unexpected end of input")
	ErrSyntaxError          = errors.New("syntax error")
	ErrInvalidUTF8          = errors.New("invalid UTF-8")
)

type Token struct {
	Type  TokenType
	Value []byte
}

type TokenType int

const (
	TOKEN_TYPE_QUERY_END TokenType = iota

	TOKEN_TYPE_SELECT

	TOKEN_TYPE_FROM
	TOKEN_TYPE_FROM_EXPRESSION

	TOKEN_TYPE_TABLE_NAME
	TOKEN_TYPE_ASTERISK
	TOKEN_TYPE_STRING
)

var (
	NullToken = Token{}
)

type State int

const (
	STATE_QUERY_START State = iota

	STATE_S
	STATE_SE
	STATE_SEL
	STATE_SELE
	STATE_SELEC
	STATE_SELECT

	STATE_SELECT_EXPRESSION
	STATE_SELECT_EXPRESSION_ASTERISK
	STATE_SELECT_EXPRESSION_STRING_LITERAL
	STATE_SELECT_POST_EXPRESSION

	STATE_F
	STATE_FR
	STATE_FRO
	STATE_FROM

	STATE_FROM_EXPRESSION
	STATE_FROM_EXPRESSION_TABLE_NAME
	STATE_FROM_POST_EXPRESSION

	STATE_FROM_TABLE

	STATE_QUERY_END
)

type Scanner struct {
	Input  []byte
	State  State
	Cursor int

	inputLen   int
	valueStart int
}

func NewScanner(input []byte) Scanner {
	return Scanner{
		Input:      input,
		State:      STATE_QUERY_START,
		Cursor:     0,
		inputLen:   len(input),
		valueStart: 0,
	}
}

func (scanner *Scanner) Next() (Token, error) {
stateLoop:
	for {
		switch scanner.State {
		case STATE_QUERY_START:
			{
				b, err := scanner.NextByteSkipWhiteSpace()
				if err != nil {
					return NullToken, err
				}

				switch b {
				case 'S', 's':
					{
						scanner.Cursor++
						scanner.State = STATE_S
						continue stateLoop
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_S:
			{
				b, err := scanner.ExpectByte()
				if err != nil {
					return NullToken, err
				}

				switch b {
				case 'E', 'e':
					{
						scanner.Cursor++
						scanner.State = STATE_SE
						continue stateLoop
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_SE:
			{
				b, err := scanner.ExpectByte()
				if err != nil {
					return NullToken, err
				}

				switch b {
				case 'L', 'l':
					{
						scanner.Cursor++
						scanner.State = STATE_SEL
						continue stateLoop
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_SEL:
			{
				b, err := scanner.ExpectByte()
				if err != nil {
					return NullToken, err
				}

				switch b {
				case 'E', 'e':
					{
						scanner.Cursor++
						scanner.State = STATE_SELE
						continue stateLoop
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_SELE:
			{
				b, err := scanner.ExpectByte()
				if err != nil {
					return NullToken, err
				}

				switch b {
				case 'C', 'c':
					{
						scanner.Cursor++
						scanner.State = STATE_SELEC
						continue stateLoop
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_SELEC:
			{
				b, err := scanner.ExpectByte()
				if err != nil {
					return NullToken, err
				}

				switch b {
				case 'T', 't':
					{
						scanner.Cursor++
						scanner.State = STATE_SELECT
						return Token{Type: TOKEN_TYPE_SELECT, Value: []byte("SELECT")}, nil
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_SELECT:
			{
				if scanner.Cursor >= scanner.inputLen {
					return NullToken, ErrUnexpectedEndOfInput
				}

				b := scanner.Input[scanner.Cursor]
				switch b {
				case ' ':
					{
						scanner.Cursor++
						scanner.State = STATE_SELECT_EXPRESSION
						continue
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_SELECT_EXPRESSION:
			{
				b, err := scanner.NextByteSkipWhiteSpace()
				if err != nil {
					return NullToken, err
				}

				switch {
				case b == '*':
					{
						scanner.Cursor++
						scanner.State = STATE_SELECT_EXPRESSION_ASTERISK
						return Token{Type: TOKEN_TYPE_ASTERISK, Value: []byte{'*'}}, nil
					}
				case b == '\'':
					{
						scanner.Cursor++
						scanner.valueStart = scanner.Cursor
						scanner.State = STATE_SELECT_EXPRESSION_STRING_LITERAL
						continue
					}
				case b >= '0' || b <= '9':
					{
						// todo
						panic("SELECT numeric literals not supported")
					}
				case (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_':
					{
						// todo
						panic("SELECT field names not supported")
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_SELECT_EXPRESSION_ASTERISK:
			{
				if scanner.Cursor >= scanner.inputLen {
					scanner.State = STATE_QUERY_END
					continue
				}

				b := scanner.Input[scanner.Cursor]
				switch b {
				case ',':
					{
						scanner.Cursor++
						scanner.State = STATE_SELECT_EXPRESSION
						continue stateLoop
					}
				case ' ':
					{
						scanner.Cursor++
						scanner.State = STATE_SELECT_POST_EXPRESSION
						continue stateLoop
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_SELECT_EXPRESSION_STRING_LITERAL:
			{
				escape := false
				for {
					if scanner.Cursor >= scanner.inputLen {
						return NullToken, ErrUnexpectedEndOfInput
					}

					b := scanner.Input[scanner.Cursor]
					switch b {
					case '\'':
						{
							if escape {
								scanner.Cursor++
								continue
							}

							v := scanner.TakeValueSlice()
							if !utf8.Valid(v) {
								return NullToken, ErrInvalidUTF8
							}

							t := Token{Type: TOKEN_TYPE_STRING, Value: v}
							scanner.Cursor++
							scanner.State = STATE_SELECT_POST_EXPRESSION
							return t, nil
						}
					case '\\':
						{
							if escape {
								escape = false
							} else {
								escape = true
							}
							scanner.Cursor++
							continue stateLoop
						}
					default:
						{
							scanner.Cursor++
							continue
						}
					}
				}
			}
		case STATE_SELECT_POST_EXPRESSION:
			{
				scanner.SkipWhitespace()

				if scanner.Cursor >= scanner.inputLen {
					scanner.State = STATE_QUERY_END
					continue
				}

				b, err := scanner.ExpectByte()
				if err != nil {
					return NullToken, err
				}

				switch b {
				case ',':
					{
						scanner.Cursor++
						scanner.State = STATE_SELECT_EXPRESSION
						continue
					}
				case 'F', 'f':
					{
						scanner.Cursor++
						scanner.State = STATE_F
						continue
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_F:
			{
				b, err := scanner.ExpectByte()
				if err != nil {
					return NullToken, err
				}

				switch b {
				case 'R', 'r':
					{
						scanner.Cursor++
						scanner.State = STATE_FR
						continue stateLoop
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_FR:
			{
				b, err := scanner.ExpectByte()
				if err != nil {
					return NullToken, err
				}

				switch b {
				case 'O', 'o':
					{
						scanner.Cursor++
						scanner.State = STATE_FRO
						continue stateLoop
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_FRO:
			{
				b, err := scanner.ExpectByte()
				if err != nil {
					return NullToken, err
				}

				switch b {
				case 'M', 'm':
					{
						scanner.Cursor++
						scanner.State = STATE_FROM
						continue stateLoop
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_FROM:
			{
				b, err := scanner.ExpectByte()
				if err != nil {
					return NullToken, err
				}

				switch b {
				case ' ':
					{
						scanner.Cursor++
						scanner.State = STATE_FROM_EXPRESSION
						continue stateLoop
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_FROM_EXPRESSION:
			{
				b, err := scanner.NextByteSkipWhiteSpace()
				if err != nil {
					return NullToken, err
				}

				switch {
				case (b >= 'A' && b <= 'Z'):
					{
						scanner.valueStart = scanner.Cursor
						scanner.Cursor++
						scanner.State = STATE_FROM_EXPRESSION_TABLE_NAME
						continue
					}
				case b == '(':
					{
						// todo
						panic("sub query not implemented")
					}
				default:
					{
						return NullToken, ErrSyntaxError
					}
				}
			}
		case STATE_FROM_EXPRESSION_TABLE_NAME:
			{
				for {
					if scanner.Cursor >= scanner.inputLen {
						v := scanner.TakeValueSlice()
						scanner.State = STATE_FROM_POST_EXPRESSION
						return Token{Type: TOKEN_TYPE_TABLE_NAME, Value: v}, nil
					}

					b := scanner.Input[scanner.Cursor]
					switch {
					case (b >= 'A' && b <= 'Z') || b == '_':
						{
							scanner.Cursor++
							continue
						}
					case b == ' ':
						{
							v := scanner.TakeValueSlice()
							t := Token{Type: TOKEN_TYPE_TABLE_NAME, Value: v}
							scanner.Cursor++
							scanner.State = STATE_FROM_POST_EXPRESSION
							return t, nil
						}
					default:
						{
							return NullToken, ErrSyntaxError
						}
					}
				}
			}
		case STATE_FROM_POST_EXPRESSION:
			{
				scanner.SkipWhitespace()

				if scanner.Cursor >= scanner.inputLen {
					scanner.State = STATE_QUERY_END
					continue
				}

				panic("post from expression not implemented")
			}
		case STATE_QUERY_END:
			{
				scanner.SkipWhitespace()

				// todo stack len check

				if scanner.Cursor >= scanner.inputLen {
					return Token{Type: TOKEN_TYPE_QUERY_END, Value: nil}, nil
				}

				return NullToken, ErrUnexpectedEndOfInput
			}
		default:
			{
				return NullToken, errors.New("unreachable")
			}
		}
	}
}

func (scanner *Scanner) NextByteSkipWhiteSpace() (byte, error) {
	scanner.SkipWhitespace()

	if scanner.Cursor >= scanner.inputLen {
		return 0, ErrUnexpectedEndOfInput
	}

	return scanner.Input[scanner.Cursor], nil
}

func (scanner *Scanner) SkipWhitespace() {
	for scanner.Cursor < scanner.inputLen {
		b := scanner.Input[scanner.Cursor]
		switch b {
		case ' ', '\t', '\r', '\n':
			{
				scanner.Cursor++
				continue
			}
		default:
			{
				return
			}
		}
	}
}

func (scanner *Scanner) TakeValueSlice() []byte {
	v := scanner.Input[scanner.valueStart:scanner.Cursor]
	scanner.valueStart = scanner.Cursor

	return v
}

func (scanner *Scanner) ExpectByte() (byte, error) {
	if scanner.Cursor >= scanner.inputLen {
		return 0, ErrUnexpectedEndOfInput
	}

	return scanner.Input[scanner.Cursor], nil
}
