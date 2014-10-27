###Javascript based action models

These actionmodels are javascript instead of json. A simple .js file for genesis doug can 
be found in the test folder, along with a main that uses it to get a value from the gendoug contract.

This is a trivial example. Action models (and other stuff too) will not live in the vm. Instead 
they would be stored in files (like action models), and when someone runs an action it would get 
the script from the file/pre-populated map, and then the vm would run the code.

