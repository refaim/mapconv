from parsimonious.grammar import Grammar
from parsimonious.nodes import NodeVisitor
from collections import OrderedDict

import re
import sys

# http://stackoverflow.com/questions/1175208/elegant-python-function-to-convert-camelcase-to-camel-case
def camel_to_underscore(name):
    s1 = re.sub('(.)([A-Z][a-z]+)', r'\1_\2', name)
    return re.sub('([a-z0-9])([A-Z])', r'\1_\2', s1).lower()


class Type(object):
    def __init__(self):
        super(Type, self).__init__()
        self.min_size = float('inf')
        self.max_size = float('inf')
        self.attrs_resolved = False


class IntegralType(Type):
    def __init__(self, size, signed):
        super(IntegralType, self).__init__()
        self.byte_size = size
        self.little_endian = True
        self.min_size = size
        self.max_size = size
        self.signed = signed
        if signed:
            self.min_value = 2 ** (size * 8 - 1) + 1
            self.max_value = 2 ** (size * 8 - 1) - 1
        else:
            self.min_value = 0
            self.max_value = 2 ** (size * 8) - 1
        self.attrs_resolved = True


class Int(object):
    def __init__(self, value):
        self.value = value
        self.min_value = value
        self.max_value = value
        self.attrs_resolved = True

    def __unicode__(self):
        return unicode(self.value)


class String(object):
    def __init__(self, min_length=0, max_length=255):
        self.min_length = min_length
        self.max_length = max_length
        self.attrs_resolved = True

    @property
    def min_size(self):
        return self.min_length + 1

    @property
    def max_size(self):
        return self.max_length + 1


class FieldRef(object):
    def __init__(self, name, level):
        self.name = name
        self.level = level
        self.ref = None

    @property
    def min_value(self):
        return self.ref.type_expr.min_value

    @property
    def max_value(self):
        return self.ref.type_expr.max_value


class BinaryOp(object):
    def __init__(self, op, left, right):
        super(BinaryOp, self).__init__()
        self.op = op
        self.left = left
        self.right = right

        if op not in ['-', '+', '*', '<<', "==", "!=", "<", "<=", ">", ">="]:
            raise NotImplementedError(op)

    @property
    def min_value(self):
        if self.op == '-':
            return self.left.min_value - self.right.max_value
        elif self.op == '+':
            return self.left.min_value + self.right.max_value
        elif self.op == '*':
            return self.left.min_value * self.right.min_value
        elif self.op == '<<':
            return (self.left.min_value << self.right.min_value) & 0xFFFFFFFF

    @property
    def max_value(self):
        if self.op == '-':
            return self.left.max_value - self.right.min_value
        elif self.op == '+':
            return self.left.max_value + self.right.min_value
        elif self.op == '*':
            return self.left.max_value * self.right.max_value
        elif self.op == '<<':
            return (self.left.max_value << self.right.max_value) & 0xFFFFFFFF


class SubType(Type):
    def __init__(self, parent_type, attributes):
        self.parent_type = parent_type
        self.attributes = attributes
        self.attrs_resolved = False

    def get_attr(self, name):
        attr = self.attributes.get(name)
        if attr is not None:
            return attr
        return getattr(self.parent_type, name)

    @property
    def min_value(self):
        return self.get_attr('min_value')

    @property
    def max_value(self):
        return self.get_attr('max_value')

    @property
    def min_length(self):
        return self.get_attr('min_length')

    @property
    def max_length(self):
        return self.get_attr('max_length')

    @property
    def little_endian(self):
        return self.get_attr('little_endian')

    @property
    def min_size(self):
        return self.parent_type.min_size

    @property
    def max_size(self):
        return self.parent_type.max_size

    @property
    def byte_size(self):
        return self.parent_type.byte_size


class FilteredType(Type):
    def __init__(self, name, parent_type):
        super(FilteredType, self).__init__()
        self.name = name
        self.parent_type = parent_type


class NamedType(object):
    def __init__(self, name, parent_type):
        self.name = name
        self.parent_type = parent_type


class Block(object):
    def __init__(self, blocks):
        super(Block, self).__init__()
        self.blocks = blocks


class Struct(Block, Type):
    def __init__(self, fields, blocks):
        super(Struct, self).__init__(blocks)
        self.fields = fields


class Array(Type):
    def __init__(self, base_type, length):
        super(Array, self).__init__()
        self.base_type = base_type
        self.length = length


class Field(object):
    def __init__(self, name, type_expr):
        self.name = name
        self.type_expr = type_expr


class IfStmt(Block):
    def __init__(self, cond, blocks):
        super(IfStmt, self).__init__(blocks)
        self.cond = cond


class TypeRef(object):
    def __init__(self, name):
        self.name = name

    def __unicode__(self):
        return self.name


class TypeUnion(Type):
    def __init__(self, expr, cases):
        super(TypeUnion, self).__init__()
        self.expr = expr
        self.cases = cases


class TypeCase(Type):
    def __init__(self, expr, type):
        super(TypeCase, self).__init__()
        self.expr = expr
        self.type = type


class Enum(Type):
    def __init__(self, base_type, items):
        super(Enum, self).__init__()
        self.base_type = base_type
        self.items = items


class Set(Type):
    def __init__(self, base_type, items):
        super(Set, self).__init__()
        self.base_type = base_type
        self.items = items


class EnumValue(object):
    def __init__(self, enum_ref, name):
        self.enum_ref = enum_ref
        self.name = name


grammar = Grammar(
r'''
body = ANY_WS? (type_decl NL ANY_WS?)+ ANY_WS?
type_decl = "type" WS IDENT WS type_expr

type_expr = array_expr / struct_expr / union_type / enum_type / type_ref / filtered_type
array_expr = "array" WS? "[" WS? expr WS? "]" WS "of" WS type_expr
type_ref = IDENT (WS? "(" WS? attribute_list? WS? ")")?
attribute_list = attribute (WS? "," WS? attribute)*
attribute = IDENT WS? "=" WS? INT

filtered_type = "$" IDENT WS? "(" WS? type_expr WS? ")"

union_type = "switch" WS expr WS? "{" (WS_WITH_NL union_item)+ WS_WITH_NL "}"
union_item = union_case / union_default
union_case = "case" WS simple_const WS? ":" WS? type_ref
union_default = "default" WS? ":" WS? type_ref

enum_type = ("enum" / "set") WS? "(" WS? type_ref WS? ")" WS? "{" (WS_WITH_NL enum_item)+ WS_WITH_NL "}"
enum_item = IDENT WS? "=" WS? expr

simple_const = enum_value / INT

expr = bin_op / atom
bin_op = atom WS? ("==" / ">" / "<<" / "-" / "+" / "*") WS? expr
atom = field_ref / enum_value / INT
field_ref = ("@" / "^"+) IDENT
enum_value = IDENT "." IDENT

struct_expr = "struct" WS "{" (((WS_WITH_NL struct_item)+ WS_WITH_NL) / ANY_WS?) "}"
struct_item = field / if_stmt
field = (IDENT / "_") WS type_expr
if_stmt = "if" WS expr WS? "{" (WS_WITH_NL struct_item)+ WS_WITH_NL "}"

IDENT = ~"[a-z][a-z0-9_]*"i
NL = ~"[\r?\n]+"
WS = ~"[\t ]+"
INT = ~"(0[xX][0-9a-fA-F]+|0|[1-9]\d*)"
ANY_WS = ~"[\r\n\t ]+"
WS_WITH_NL = ~"[\r\t ]*\n[\r\n\t ]*"
''')

class Compiler(NodeVisitor):
    def __init__(self):
        super(NodeVisitor, self).__init__()
        self.line = 1

    def visit_body(self, node, (_1, decls, _2)):
        return [decl for decl, _1, _2 in decls]

    def visit_type_decl(self, node, (_1, _2, ident, _3, type_expr)):
        return NamedType(ident, type_expr)

    def visit_struct_expr(self, node, (_1, _2, _3, _4, _5)):
        if not _4[0] or _4[0][0] == None:
            return Struct([], [])
        return Struct([], [field for _1, field in _4[0][0]])

    def visit_struct_item(self, node, (item,)):
        return item

    def visit_field(self, node, (name, _, type_expr)):
        return Field(name[0], type_expr)

    def visit_if_stmt(self, node, (_1, _2, expr, _3, _4, values, _5, _6)):
        return IfStmt(expr, [value[1] for value in values])

# struct_item = field / if_stmt
# field = (IDENT / "_") WS type_expr
# if_stmt = "if" WS expr WS? "{" (WS_WITH_NL struct_item)+ WS_WITH_NL "}"

    def visit_type_expr(self, node, expr):
        if isinstance(expr[0], (str, unicode)):
            return TypeRef(expr[0])
        return expr[0]

    def visit_type_ref(self, node, (name, attrs)):
        ref = TypeRef(name)
        if len(attrs) > 0:
            return SubType(ref, attrs[0][3][0])
        return ref

    def visit_enum_type(self, node, (_1, _2, _3, _4, base_type, _5, _6, _7, _8, items, _9, _10)):
        values = OrderedDict(item[1] for item in items)
        if _1[0] == 'set':
            return Set(base_type, values)
        return Enum(base_type, values)

    def visit_enum_item(self, node, (ident, _1, _2, _3, expr)):
        return ident, expr

    def visit_attribute_list(self, node, children):
        attributes = [children[0]]
        if len(children) > 1:
            attributes += [child[3] for child in children[1]]
        return OrderedDict(attributes)

    def visit_attribute(self, node, (name, _1, _2, _3, value)):
        return camel_to_underscore(name), value.value

    def visit_array_expr(self, node, (_1, _2, _3, _4, expr, _5, _6, _7, _8, _9, base_type)):
        return Array(base_type, expr)

    def visit_expr(self, node, children):
        return children[0]

    def visit_atom(self, node, children):
        return children[0]

    def visit_bin_op(self, node, (left, _1, op, _2, right)):
        return BinaryOp(op[0], left, right)

    def visit_field_ref(self, node, (_1, name)):
        level = -1
        if isinstance(_1[0], list):
            level = -len(_1[0]) - 1
        return FieldRef(name, level)

    def visit_filtered_type(self, node, (_1, name, _2, _3, _4, parent_type, _5, _6)):
        return FilteredType(name, parent_type)

    def visit_union_type(self, node, (_1, _2, expr, _3, _4, values, _5, _6)):
        return TypeUnion(expr, [value[1] for value in values])

    def visit_union_item(self, node, (value,)):
        return value

    def visit_union_case(self, node, (_1, _2, expr, _3, _4, _5, type)):
        return TypeCase(expr, type)

    def visit_union_default(self, node, (_1, _2, _3, _4, type)):
        return TypeCase(None, type)

    def visit_simple_const(self, node, (value,)):
        return value

    def visit_enum_value(self, node, (enum, _1, name)):
        return EnumValue(enum, name)

    def visit_IDENT(self, node, children):
        return node.text

    def visit_NL(self, node, children):
        self.line += node.text.count('\n')

    def visit_ANY_WS(self, node, children):
        self.line += node.text.count('\n')

    def visit_WS(self, node, children):
        return None

    def visit_INT(self, node, children):
        return Int(int(node.text, 0))

    def visit_(self, node, children):
        if len(children) == 0:
            return node.text
        return children

    def generic_visit(self, node, children):
        return children


def add_default_types(types):
    types['uint8'] = IntegralType(1, signed=False)
    types['uint16'] = IntegralType(2, signed=False)
    types['uint32'] = IntegralType(4, signed=False)
    types['int8'] = IntegralType(1, signed=True)
    types['int16'] = IntegralType(2, signed=True)
    types['int32'] = IntegralType(4, signed=True)
    types['cstring'] = String()


class UnresolvedType(object):
    UNKNOWN = 0
    RESOLVING = 1
    RESOLVED = 2

    def __init__(self, type):
        self.type = type
        self.status = self.UNKNOWN



def resolve_types(type_list):
    types = OrderedDict()
    reverse_types = {}
    add_default_types(types)
    default_types = types.items()

    main_type = None
    for type in type_list:
        if isinstance(type, NamedType):
            if isinstance(type.parent_type, Struct) and main_type is None:
                main_type = type.parent_type
            types[type.name] = UnresolvedType(type.parent_type)
        else:
            raise NotImplementedError

    resolved_types = []
    def resolve_expr(expr, context):
        if isinstance(expr, FieldRef):
            if len(context) < -expr.level:
                raise Exception('Invalid reference {0}{1}: struct not found'.format("@" if expr.level == -1 else "^" * (-expr.level - 1), expr.name))
            env = context[expr.level]
            field = env['fields'].get(expr.name)
            if field is None:
                raise Exception('Invalid reference {0}{1}: unknown field'.format("@" if expr.level == -1 else "^" * (-expr.level - 1), expr.name))
            expr.ref = field
        elif isinstance(expr, BinaryOp):
            expr.left = resolve_expr(expr.left, context)
            expr.right = resolve_expr(expr.right, context)
        elif isinstance(expr, Int):
            pass
        elif isinstance(expr, EnumValue):
            expr.enum_ref = types.get(expr.enum_ref)
            if expr.enum_ref is None or not isinstance(expr.enum_ref, (Enum, Set)) or expr.name not in expr.enum_ref.items:
                raise Exception("Invalid enum ref {0}::{1}".format(expr.enum_ref, expr.name))
        else:
            raise NotImplementedError('Unsupported expression type {0}'.format(expr.__class__.__name__))
        return expr

    def resolve_type(expr, context):
        if isinstance(expr, (IntegralType, String)):
            return expr
        elif isinstance(expr, UnresolvedType):
            if expr.status == UnresolvedType.RESOLVED:
                return expr.type
            elif expr.status == UnresolvedType.RESOLVING:
                raise Exception('Cycle detected')

            expr.status = UnresolvedType.RESOLVING
            expr.type = resolve_type(expr.type, context)
            expr.status = UnresolvedType.RESOLVED

            resolved_types.append(expr.type)
            return expr.type
        elif isinstance(expr, Struct):
            processed_fields = OrderedDict()
            unknown_counter = [1]
            context.append({'fields': processed_fields, 'struct': expr})

            def resolve_block(block):
                if isinstance(block, Field):
                    block.type_expr = resolve_type(block.type_expr, context)
                    if block.name == "_":
                        block.name = "Unknown{0}".format(unknown_counter[0])
                        unknown_counter[0] += 1
                    processed_fields[block.name] = block
                elif isinstance(block, IfStmt):
                    block.cond = resolve_expr(block.cond, context)
                    for subblock in block.blocks:
                        resolve_block(subblock)
                else:
                    raise NotImplementedError(expr.__class__.__name__)

            for block in expr.blocks:
                resolve_block(block)

            context.pop()
            expr.fields = processed_fields.values()
            return expr
        elif isinstance(expr, Array):
            expr.base_type = resolve_type(expr.base_type, context)
            expr.length = resolve_expr(expr.length, context)
            return expr
        elif isinstance(expr, TypeRef):
            ref_type = types.get(expr.name)
            if ref_type is None:
                raise Exception('Unsupported type {0}'.format(expr.name))
            return resolve_type(ref_type, context)
        elif isinstance(expr, SubType):
            expr.parent_type = resolve_type(expr.parent_type, context)
            if isinstance(expr.parent_type, String):
                return String(**expr.attributes)
            return expr
        elif isinstance(expr, TypeUnion):
            expr.expr = resolve_expr(expr.expr, context)
            expr.cases = [resolve_type(case, context) for case in expr.cases]
            return expr
        elif isinstance(expr, TypeCase):
            if expr.expr is not None:
                expr.expr = resolve_expr(expr.expr, context)
            expr.type = resolve_type(expr.type, context)
            return expr
        elif isinstance(expr, (Enum, Set)):
            expr.base_type = resolve_type(expr.base_type, context)
            if not isinstance(expr.base_type, IntegralType):
                raise Exception("Expected integer type")
            for name, item in expr.items.iteritems():
                expr.items[name] = resolve_expr(item, context)
            return expr
        elif isinstance(expr, FilteredType):
            expr.parent_type = resolve_type(expr.parent_type, context)
            return expr
        else:
            raise NotImplementedError(expr.__class__.__name__)

    for name, expr in types.iteritems():
        if isinstance(expr, UnresolvedType) and isinstance(expr.type, (Enum, Set)):
            types[name] = resolve_type(expr, [])
            reverse_types[types[name]] = name

    for name, type in types.iteritems():
        types[name] = resolve_type(type, [])
        reverse_types[types[name]] = name

    return main_type, OrderedDict(default_types + zip(map(reverse_types.__getitem__, resolved_types), resolved_types))


def resolve_attributes(value):
    if value.attrs_resolved:
        return

    if isinstance(value, Struct):
        for field in value.fields:
            resolve_attributes(field.type_expr)
        value.min_size = sum(field.type_expr.min_size for field in value.fields)
        value.max_size = sum(field.type_expr.max_size for field in value.fields)
    elif isinstance(value, Array):
        value.min_size = value.length.min_value * value.base_type.min_size
        value.max_size = value.length.max_value * value.base_type.max_size
    elif isinstance(value, SubType):
        resolve_attributes(value.parent_type)
    elif isinstance(value, TypeCase):
        resolve_attributes(value.type)
        value.min_size = value.type.min_size
        value.max_size = value.type.max_size
    elif isinstance(value, TypeUnion):
        for case in value.cases:
            resolve_attributes(case)
        value.min_size = min(case.min_size for case in value.cases)
        value.max_size = max(case.max_size for case in value.cases)
    elif isinstance(value, FilteredType):
        resolve_attributes(value.parent_type)
    elif isinstance(value, (Enum, Set)):
        resolve_attributes(value.base_type)
        value.min_size = value.base_type.min_size
        value.max_size = value.base_type.max_size
    else:
        raise Exception('Unsupported object {0}'.format(value.__class__.__name__))
    value.attrs_resolved = True


def get_parent_type(type):
    if isinstance(type, SubType):
        return type.parent_type
    return type

def add_indent(content, level):
    return '\n'.join(["\t" * level + line for line in content.split("\n")])

class GenerationImpossible(Exception): pass

def convert_to_go(types, package_name):
    def out(format, *args):
        result.append(format.format(*args).replace('\t', indent))

    def is_ptr_array(type):
        return (isinstance(type, Array)
            and not isinstance(type.length, Int)
            and not isinstance(type.base_type, (IntegralType, String)))

    def get_go_type(type):
        if isinstance(type, Set):
            return get_go_type(type.base_type)
        elif isinstance(type, String):
            return "string"
        elif isinstance(type, FilteredType):
            return get_go_type(type.parent_type)
        elif isinstance(type, TypeUnion):
            return "interface{}"
        elif type in type_names:
            return type_names[type]
        elif isinstance(type, SubType):
            return get_go_type(type.parent_type)
        elif isinstance(type, Array):
            p = '*' if is_ptr_array(type) else ''
            if isinstance(type.length, Int):
                return '[{}]{}{}'.format(type.length.value, p, get_go_type(type.base_type))
            else:
                return '[]{}{}'.format(p, get_go_type(type.base_type))
        return type.__class__.__name__

    def format_go_expr(expr, context):
        def resolve_ref(r):
            try:
                return "int({0}.{1})".format(context[r.level], type_names[r.ref])
            except IndexError:
                raise GenerationImpossible()
        if isinstance(expr, Int):
            return unicode(expr.value)
        elif isinstance(expr, BinaryOp):
            return u"({0} {1} {2})".format(format_go_expr(expr.left, context), expr.op, format_go_expr(expr.right, context))
        elif isinstance(expr, FieldRef):
            return resolve_ref(expr)
        elif isinstance(expr, EnumValue):
            return "int({0}{1})".format(type_names[expr.enum_ref], expr.name)
        raise NotImplementedError(expr.__class__.__name__)

    type_names = dict((type, name) for name, type in types.iteritems())
    result = []
    indent = ' ' * 4

    out('package {}\n\n'.format(package_name))
    out('type Serializer interface {{\n')
    out('\tUInt8(*uint8)\n')
    out('\tUInt16(*uint16)\n')
    out('\tUInt32(*uint32)\n')
    out('\tInt8(*int8)\n')
    out('\tInt16(*int16)\n')
    out('\tInt32(*int32)\n')
    out('\tCString(*string)\n')
    out('\tIsReader() bool\n')
    out('}}\n\n')

    for name, type in types.iteritems():
        if isinstance(type, Enum):
            out('type {0} {1}\n', name, get_go_type(type.base_type))

    out('\n')

    for name, type in types.iteritems():
        if isinstance(type, Struct):
            out('type {0} struct {{\n', name, type_names[type])
            if type.fields:
                max_len = max(len(field.name) for field in type.fields)
                for field in type.fields:
                    type_names[field] = field.name
                    type_expr = get_go_type(field.type_expr)
                    out('\t{0} {1}\n', field.name.ljust(max_len), type_expr)
            out('}}\n\n')

    for name, type in types.iteritems():
        if isinstance(type, Enum):
            out('const (\n')
            first = True

            max_len = max(len(name) for name in type.items)
            for enum_name, expr in type.items.iteritems():
                if first:
                    out('\t{1}{0} {1} = {2}\n', enum_name.ljust(max_len), name, format_go_expr(expr, []))
                    first = False
                else:
                    out('\t{1}{0} = {2}\n', enum_name.ljust(max_len + len(name) + 1), name, format_go_expr(expr, []))
            out(')\n\n')
        elif isinstance(type, Set):
            out('const (\n')
            max_len = max(len(name) for name in type.items)
            for enum_name, expr in type.items.iteritems():
                out('\t{1}{0} = {2}\n', enum_name.ljust(max_len), name, format_go_expr(expr, []))
            out(')\n\n')

    counter = [0]

    def generate_blocks(blocks, path, context, depth):
        indent = '\t' * depth
        for block in blocks:
            if isinstance(block, Field):
                field_path = '(&({}.{}))'.format(path, block.name)
                generate_reader(block.type_expr, field_path, context, depth)
            elif isinstance(block, IfStmt):
                out('{}if {} {{\n', indent, format_go_expr(block.cond, context))
                generate_blocks(block.blocks, path, context, depth + 1)
                out('{}}}\n', indent)
            else:
                raise Exception("Unsupported block type {0}".format(block.__class__.__name__))

    def generate_reader(type, path, context, depth):
        indent = '\t' * depth
        if isinstance(type, Struct):
            new_context = context + [path]
            generate_blocks(type.blocks, path, new_context, depth)
        elif isinstance(type, IntegralType):
            assert type.little_endian == True
            out('{}b.{}Int{}({})\n'.format(indent, '' if type.signed else 'U', type.byte_size * 8, path))
        elif isinstance(type, String):
            # TODO min max length
            out('{}b.CString({})\n'.format(indent, path))
        elif isinstance(type, Set):
            generate_reader(type.base_type, path, context, depth)
        elif isinstance(type, Enum):
            generate_reader(type.base_type, '(*{})({})'.format(get_go_type(type.base_type), path), context, depth)
        elif isinstance(type, FilteredType):
            generate_reader(type.parent_type, path, context, depth)
        elif isinstance(type, TypeUnion):
            out('{}switch {} {{\n', indent, format_go_expr(type.expr, context))
            for case in type.cases:
                counter[0] += 1
                if case.expr is None:
                    out('{}default:\n', indent)
                else:
                    out('{}case {}:\n', indent, format_go_expr(case.expr, context))
                var = 'v{}'.format(counter[0])
                out('\t{}if b.IsReader() {{\n', indent)
                out('\t\t{}var {} {}\n', indent, var, get_go_type(case.type))
                generate_reader(case.type, var, context, depth + 2)
                out('\t\t{}*{} = &{}\n', indent, path, var)
                out('\t{}}} else {{\n', indent)
                out('\t\t{}{} := (*{}).(*{})\n', indent, var, path, get_go_type(case.type))
                generate_reader(case.type, var, context, depth + 2)
                out('\t{}}}\n', indent)
            out('{}}}\n', indent)
        elif isinstance(type, Array):
            counter[0] += 1
            size = 'size{}'.format(counter[0])
            i = 'i{}'.format(counter[0])
            v = 'v{}'.format(counter[0])
            dynamic = not isinstance(type.length, Int)

            new_path = '({}(*{})[{}])'.format('' if is_ptr_array(type) else '&', path, i)

            out('{}{} := {}\n'.format(indent, size, format_go_expr(type.length, context)))

            if dynamic:
                out('{}if b.IsReader() {{\n', indent)
                out('\t{}*{} = make({}, {})\n', indent, path, get_go_type(type), size)
                if not isinstance(type.base_type, IntegralType):
                    out('\t{0}for {1} := (0); {1} < {2}; {1}++ {{\n', indent, i, size)
                    out('\t\t{0}{3}{1} = &{2}{{}}\n', indent, new_path, get_go_type(type.base_type), '' if is_ptr_array(type) else '*')
                    out('\t{}}}\n', indent)
                out('{}}}\n', indent)

            out('{0}for {1} := (0); {1} < {2}; {1}++ {{\n', indent, i, size)
            generate_reader(type.base_type, new_path, context, depth + 1)
            out('{}}}\n', indent)
        elif isinstance(type, SubType):
            generate_reader(type.parent_type, path, context, depth)
        else:
            raise Exception("Unsupported type {}".format(type))

    for name, type in types.iteritems():
        if isinstance(type, Struct):
            old_len = len(result)
            try:
                out('func (s *{}) Serialize(b Serializer) {{\n', name)
                generate_reader(type, "s", [], 1)
                out('}}\n\n')
            except GenerationImpossible:
                print get_go_type(type)
                result = result[:old_len]

    return ''.join(result)


input_file = sys.argv[1]
package_name = sys.argv[2]
output_file = sys.argv[3]

text = open(input_file).read().decode('utf-8')
text = re.sub(r'//.*($|[\r\n])', '\\1', text)

nodes = grammar.parse(text + '\n')
if nodes is not None:
    types = Compiler().visit(nodes)
    main_type, types = resolve_types(types)

    type_names = {}
    for name, type in types.iteritems():
        resolve_attributes(type)
        type_names[type] = name

    out = convert_to_go(types, package_name)
    with open(output_file, 'w') as of:
        of.write(out)

    # template = open('test.template.cpp').read()
    # template = template.replace('{STRUCTS}', structs)
    # template = template.replace('{FSM}', fsm)
    # template = template.replace('{PRINTER}', printer)
    # template = template.replace('{FIELDS}', fields)

else:
    print 'Parse failed'
    sys.exit(1)