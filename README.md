The 9spell command runs the plan9port "9 spell" program on each
of the files supplied as arguments. It prints output like:

	file0:44+/teh/
	file0:63+/frgo/
	file1:0+/fner/

A program such as acme can read those addresses and navigate
to the misspelled word.

If a filename ends in ".tex", that file is piped through the plan9port
"9 detex" program before "9 spell".
