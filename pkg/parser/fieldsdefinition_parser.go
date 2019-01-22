package parser

import (
	"github.com/jensneuse/graphql-go-tools/pkg/document"
	"github.com/jensneuse/graphql-go-tools/pkg/lexing/keyword"
	"github.com/jensneuse/graphql-go-tools/pkg/lexing/position"
)

func (p *Parser) parseFieldsDefinition(index *[]int) (err error) {

	if hasOpen := p.peekExpect(keyword.CURLYBRACKETOPEN, true); !hasOpen {
		return
	}

	var description *document.ByteSliceReference
	var startPosition *position.Position

	for {
		next := p.l.Peek(true)

		switch next {
		case keyword.STRING:
			stringToken := p.l.Read()
			description = &stringToken.Literal
			startPosition = &stringToken.TextPosition
		case keyword.CURLYBRACKETCLOSE:
			p.l.Read()
			return nil
		case keyword.IDENT, keyword.TYPE:

			fieldIdent := p.l.Read()
			definition := p.makeFieldDefinition()

			if description != nil {
				definition.Description = *description
				description = nil
			}

			if startPosition != nil {
				definition.Position.MergeStartIntoStart(*startPosition)
				startPosition = nil
			} else {
				definition.Position.MergeStartIntoStart(fieldIdent.TextPosition)
			}

			definition.Name = fieldIdent.Literal

			err = p.parseArgumentsDefinition(&definition.ArgumentsDefinition)
			if err != nil {
				return err
			}

			_, err = p.readExpect(keyword.COLON, "parseFieldsDefinition")
			if err != nil {
				return err
			}

			err = p.parseType(&definition.Type)
			if err != nil {
				return err
			}

			definition.Position.MergeStartIntoEnd(p.TextPosition())

			err = p.parseDirectives(&definition.Directives)
			if err != nil {
				return err
			}

			*index = append(*index, p.putFieldDefinition(definition))

		default:
			invalid := p.l.Read()
			return newErrInvalidType(invalid.TextPosition, "parseFieldsDefinition", "string/curly bracket close/ident", invalid.Keyword.String())
		}
	}
}
