# DocDork

#### A lightweight Bing document finder written in Go.

## __Description__

A lot of data can be found through search engine dorks and documents are no exception. DocDork uses a relatively lightweight method to dork bing and find documents exposed to the internet which can potentially reveal harmful information. 

This project can help an organization determine the data that they may be exposing inadvertently. 

```
	
      __ \                __ \                |    
      |   |   _ \    __|  |   |   _ \    __|  |  / 
      |   |  (   |  (     |   |  (   |  |       <     
      ___/  \___/  \___| ____/  \___/  _|    _|\_\

	      DocDork (v0.1)
		


Usage of DocDork.exe:
  -d string
        Specify the domain to search.
  -e string
        Specify the extension type. No argument will query every type. Options: docx, pptx, xlsx, pdf
  -m
        Used to output metadata. (Only works with open XML extensions)
```

### __Features__

- Discovers documents exposed to the internet.

- Can be tailored to target specific filetypes. 

- For open XML documents (I.E. docx, pptx, and xlsx), DocDork can be configured to output metadata.

### __Limitations and Considerations__

- Only yields the first 10 documents per filetype (Bing results per page).

- Metadata is not configured for PDFs.

## __Basic Usage__

DocDork can be compiled to run from either Linux or Windows, thanks to Go. 

After compilation, DocDork can be ran by simply specifying a domain. 

```
$ DocDork -d example.com
[...]
Domain: *.example.com
[Info] Seaching for docx extensions.
[Doc Found] https://www.example.com/protected/secret.docx
```

Additionally, DocDork can also accept specific document filetypes. If no type is specified, it will query for all four supported extensions. DocDork will return the metadata for open XML if the `-m` switch is enabled.

```
$ DocDork -d example.com -e docx -m
[...]
Domain: *.example.com
[Info] Seaching for docx extensions.
[Doc Found] https://www.example.com/protected/secret.docx
[Metadata]
  Creator               : Doe, John
  Last Modified By      : Jane, Mary
  Company               : Example Co.
  Application           : Microsoft Office Word
  Version               : 2016
```

---

### __Versions__

__0.0.1:__

- Initial release

### __Disclaimer__

This open source project is meant to be used with explicit authorization from any entity it affects (perceived or actual). This programs use in conjunction with offensive security tools should only take place in an approved assessment of an organization's security or for authorized research. Misuse of this software is not the responsibility of the author.
