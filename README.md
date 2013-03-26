The 9spell command runs the plan9port "9 spell" program on each
of the files supplied as arguments. It prints output like:

	file_a:44+/teh/
	file_a:63+/frgo/
	file_b:0+/fner/

A program such as acme can read those addresses and navigate
to the misspelled word.

If a filename ends in ".tex", that file is piped through the plan9port
"9 detex" program before "9 spell".
