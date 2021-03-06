package phpserialize

import "strconv"

const (
	TOKEN_NULL              rune = 'N'
	TOKEN_BOOL              rune = 'b'
	TOKEN_INT               rune = 'i'
	TOKEN_FLOAT             rune = 'd'
	TOKEN_STRING            rune = 's'
	TOKEN_ARRAY             rune = 'a'
	TOKEN_OBJECT            rune = 'O'
	TOKEN_OBJECT_SERIALIZED rune = 'C'
	TOKEN_REFERENCE         rune = 'R'
	TOKEN_REFERENCE_OBJECT  rune = 'r'
	TOKEN_SPL_ARRAY         rune = 'x'
	TOKEN_SPL_ARRAY_MEMBERS rune = 'm'

	SEPARATOR_VALUE_TYPE rune = ':'
	SEPARATOR_VALUES     rune = ';'

	DELIMITER_STRING_LEFT  rune = '"'
	DELIMITER_STRING_RIGHT rune = '"'
	DELIMITER_OBJECT_LEFT  rune = '{'
	DELIMITER_OBJECT_RIGHT rune = '}'

	FORMATTER_FLOAT     byte = 'g'
	FORMATTER_PRECISION int  = 17
)

var (
	debugMode = false
)

func Debug(value bool) {
	debugMode = value
}

func NewPhpObject(className string) *PhpObject {
	return &PhpObject{
		className: className,
		members:   PhpArray{},
	}
}

type SerializedDecodeFunc func(string) (PhpValue, error)

type SerializedEncodeFunc func(PhpValue) (string, error)

type PhpValue interface{}

type PhpArray map[PhpValue]PhpValue

type PhpSlice []PhpValue

type PhpObject struct {
	className string
	members   PhpArray
}

func (self *PhpObject) GetClassName() string {
	return self.className
}

func (self *PhpObject) SetClassName(name string) *PhpObject {
	self.className = name
	return self
}

func (self *PhpObject) GetMembers() PhpArray {
	return self.members
}

func (self *PhpObject) SetMembers(members PhpArray) *PhpObject {
	self.members = members
	return self
}

func (self *PhpObject) GetPrivate(name string) (v PhpValue, ok bool) {
	v, ok = self.members["\x00"+self.className+"\x00"+name]
	return
}

func (self *PhpObject) SetPrivate(name string, value PhpValue) *PhpObject {
	self.members["\x00"+self.className+"\x00"+name] = value
	return self
}

func (self *PhpObject) GetProtected(name string) (v PhpValue, ok bool) {
	v, ok = self.members["\x00*\x00"+name]
	return
}

func (self *PhpObject) SetProtected(name string, value PhpValue) *PhpObject {
	self.members["\x00*\x00"+name] = value
	return self
}

func (self *PhpObject) GetPublic(name string) (v PhpValue, ok bool) {
	v, ok = self.members[name]
	return
}

func (self *PhpObject) SetPublic(name string, value PhpValue) *PhpObject {
	self.members[name] = value
	return self
}

func NewPhpObjectSerialized(className string) *PhpObjectSerialized {
	return &PhpObjectSerialized{
		className: className,
	}
}

type PhpObjectSerialized struct {
	className string
	data      string
	value     PhpValue
}

func (self *PhpObjectSerialized) GetClassName() string {
	return self.className
}

func (self *PhpObjectSerialized) SetClassName(name string) *PhpObjectSerialized {
	self.className = name
	return self
}

func (self *PhpObjectSerialized) GetData() string {
	return self.data
}

func (self *PhpObjectSerialized) SetData(data string) *PhpObjectSerialized {
	self.data = data
	return self
}

func (self *PhpObjectSerialized) GetValue() PhpValue {
	return self.value
}

func (self *PhpObjectSerialized) SetValue(value PhpValue) *PhpObjectSerialized {
	self.value = value
	return self
}

func NewPhpSplArray(array, properties PhpValue) *PhpSplArray {
	if array == nil {
		array = make(PhpArray)
	}

	if properties == nil {
		properties = make(PhpArray)
	}

	return &PhpSplArray{
		array:      array,
		properties: properties,
	}
}

type PhpSplArray struct {
	flags      int
	array      PhpValue
	properties PhpValue
}

func (self *PhpSplArray) GetFlags() int {
	return self.flags
}

func (self *PhpSplArray) SetFlags(value int) {
	self.flags = value
}

func (self *PhpSplArray) GetArray() PhpValue {
	return self.array
}

func (self *PhpSplArray) SetArray(value PhpValue) {
	self.array = value
}

func (self *PhpSplArray) GetProperties() PhpValue {
	return self.properties
}

func (self *PhpSplArray) SetProperties(value PhpValue) {
	self.properties = value
}

func PhpValueString(p PhpValue) (res string) {
	res, _ = p.(string)
	return
}

func PhpValueBool(p PhpValue) (res bool) {
	switch p.(type) {
	case bool:
		res, _ = p.(bool)
	case string:
		str, _ := p.(string)
		res, _ = strconv.ParseBool(str)
	}
	return
}

func PhpValueInt(p PhpValue) (res int) {
	switch p.(type) {
	case int:
		res, _ = p.(int)
	case int8:
		intVal, _ := p.(int8)
		res = int(intVal)
	case int16:
		intVal, _ := p.(int16)
		res = int(intVal)
	case int32:
		intVal, _ := p.(int32)
		res = int(intVal)
	case int64:
		intVal, _ := p.(int64)
		res = int(intVal)
	case uint:
		intVal, _ := p.(uint)
		res = int(intVal)
	case uint8:
		intVal, _ := p.(uint8)
		res = int(intVal)
	case uint16:
		intVal, _ := p.(uint16)
		res = int(intVal)
	case uint32:
		intVal, _ := p.(uint32)
		res = int(intVal)
	case uint64:
		intVal, _ := p.(uint64)
		res = int(intVal)
	case string:
		str, _ := p.(string)
		res, _ = strconv.Atoi(str)
	}
	return
}

func PhpValueInt64(p PhpValue) (res int64) {
	switch p.(type) {
	case int64:
		res = p.(int64)
	default:
		res = int64(PhpValueInt(p))
	}
	return
}

func PhpValueUInt(p PhpValue) (res uint) {
	switch p.(type) {
	case uint:
		res = p.(uint)
	default:
		res = uint(PhpValueInt(p))
	}
	return
}

func PhpValueUInt64(p PhpValue) (res uint64) {
	switch p.(type) {
	case uint64:
		res = p.(uint64)
	default:
		res = uint64(PhpValueInt(p))
	}
	return
}

func PhpValueFloat64(p PhpValue) (res float64) {
	switch p.(type) {
	case float64:
		res, _ = p.(float64)
	case string:
		str, _ := p.(string)
		res, _ = strconv.ParseFloat(str, 64)
	default:
		return float64(PhpValueInt(p))
	}
	return
}
