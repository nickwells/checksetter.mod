/*
Package checksetter provides a generic Setter that can be used to construct
lists of check functions. It also provides several pre-constructed Parsers
for converting a string into a slice of check.ValCk functions of the
appropriate type. To construct the Setter you should set the Value as with
other setters but you will also need to set the Parser to use. You can
retrieve a Parser value with which to initialise the Setter with the
FindParser func which you will need to call with the type set to the type of
value to be checked and the checker name set to one of the const
...CheckerName values.

If you choose to write your own Parser you should do it by calling the
MakeParser func which will register the Parser so that it can be retrieved
with the FindParser func. This will then also allow the Setter to provide
correct AllowedValues.
*/
package checksetter
