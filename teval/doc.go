/*

Package teval provides means for evaluating tables
of the following format:

  | - separator
  r - row number (positive integer)
  h - header (string)
  e - number or expression (= arg1 *operand* arg2), must not be empty

    |h|h|...
  r |e|e|...
  ...

Note that first row defines header (column names) and it's
first value is missing, indicating empty top left cell of the table.

Row numbers are not necessarily in sequential and in ascending order.

Example:

  ,A,B,Cell
  1,1,0,1
  2,2,=A1+Cell30,0
  30,0,=B1+A1,5

  represents the following table

      A  B           Cell
  1   1  0           1
  2   2  =A1+Cell30  0
  30  0  =B1+A1      5

*/
package teval
