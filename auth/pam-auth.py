from simplepam import authenticate
import sys
print 'OK' if authenticate(sys.argv[1], sys.argv[2]) else 'ERR'
