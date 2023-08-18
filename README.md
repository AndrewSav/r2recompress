# Recompress Remnant II save files

Usage: r2recompress.exe [options] inputFile outputFile

Specify one of the three modes:
- `-d` - decompress: input file is profile.sav or save_x.sav, output file - decompressed file
- `-c` - compress: input file is a decompressed sav file, output file is a game save file that the game will accept (if it's a legit decompressed save otherwise)
- `-s` - dump strings: input file is a decompressed sav file, output file a plain text file you can open in notepad with all strings form the save dumped

Other options:
- `-q` -  decompress only - only print the first warning if there are warnings
- `-qq` - decompress only - do not print warnings
- `-v` - compress/decompress only - print debug data about save file structure

If the save file does not match assumptions, a warning or warnings will be displayed. If there are any the chance that the file could be recompressed to a working game save is severely reduced. This programm will probably need to be fixed if you ever see a warning
