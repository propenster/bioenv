# bioenv
BioEnv is a virtual environment for bioinformatics exploration, analysis and research

## Features
- Virtual environment for bioinformatics workflow
- Package/dependency management for bioinformatics workflow


### How to use
- After having installed and set up bioenv on your machine

**Create/Initialize a new bio env**
```bash
bioenv init . myfirstenv
# note . is the specified working directory where you want to setup an env

```
This command generates a new bioinformatics virtual environment and creates a bioenv.json configuration file in the specified work directory.

**Install a new tool**
```bash
bioenv install gatk

```
This command installs a new tool into your environment and modifies the config bioenv.json accordingly 

**Call/Use any tool in your environment**
Still thinking about how to improve this

```bash
bioenv call gatk --java-options "-Xmx4G" [program arguments]

```
or 

```bash

bioenv call bwa index ref.fasta
bioenv call samtools faidx ref.fasta

```

**Stop/Quit virtual environment**

```bash

bioenv quit

```

**Export virtual environment so you can be able to share it with your colleagues**

```bash
bioenv export

```

