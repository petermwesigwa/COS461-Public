#!/bin/bash
./longstring $1 > input.txt
./client-c 127.0.0.1 25000 < input.txt > sockout.txt
echo "============================================"
echo "Input to client vs output from server:"
diff input.txt output.txt
echo ""
echo "============================================"
echo "Input to client vs client to server:"
diff input.txt sockout.txt
echo ""
echo "============================================"
echo "Client to server vs output from server"
diff sockout.txt output.txt
echo ""
rm input.txt output.txt sockout.txt