# CVital
The goal of this project is create an online resource for job seekers looking to improve their CVs. There are two main interlinked functions:

1) Provide a tool to allow for the creation of downloadable CV pdfs. Eventually the software should provide assistance to users through stylish templates and spell checking
2) Keep the users CV details on file, which are then searchable by companies looking for employees. The companies create separate profiles that allow them to search through the created profiles.

The goal is to try and put the job seekers first. This will be made available on a website, but ideally a mobile ready one that can be accessed by people without access to a PC.

## Running
To run this program:

To start running system dependencies, install docker and then:
```bash
docker-compose up -d
```

To run the backend, first make a copy of config_example.yml rename it to config.yml. Then run:
```bash
go run main.go
```

To run the frontend (this is highly rudimentary):
```bash
cd app
npm start
```