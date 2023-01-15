VERSION 0.6
FROM mcr.microsoft.com/dotnet/sdk:7.0
WORKDIR /src

ci:
  BUILD +test

cd:
  BUILD +dockerize-web-application

deps:
  # consider copying only csprojs so you can cache the restore
  # COPY src\Template.Web.Application\Template.Web.Application.csproj src\Template.Web.Application
 
  COPY src src
  RUN dotnet restore src/Template-Solution.sln

build:
  FROM +deps
  COPY src src

  RUN dotnet build --no-restore src/Template-Solution.sln

test:
  COPY scripts/test.sh .

  FROM +build
  COPY src src

  RUN test.sh

  RUN mkdir -p test-results
  RUN cp src/**/*.trx test-results

  SAVE ARTIFACT test-results AS LOCAL test-results

publish-web-application:
  FROM +build
  COPY src src

  RUN dotnet publish --no-restore --no-build src/Template.Web.Application/Template.Web.Application.csproj -o Template.Web.Application

  SAVE ARTIFACT Template.Web.Application

dockerize-web-application:
  FROM mcr.microsoft.com/dotnet/aspnet:7.0
  COPY +publish-web-application/Template.Web.Application .
  ENTRYPOINT ["dotnet", "Template.Web.Application.dll"]
  SAVE IMAGE --push dotnet-template/template:web-application
