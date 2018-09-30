go build -o Game ./ 

# This gives information about elements that could be simplified if built with the args. (Thanks, acln on Discord!)
# go build -gcflags="-m -m" -o Game ./ 

./Game
