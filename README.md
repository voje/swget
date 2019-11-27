# swget
Smart wget.  

Some more golang practise.  

## Command structure
```bash
$ swget url [file [--exact|--highest-version]]
```
In case of [file]:
By default, find all files that match the filename and present an interactive file selection.  
`--exact` Find exact filename match or return error exit code.  
`--highest-version` Return file with the highest version (regex for three-number-version in file name).  

