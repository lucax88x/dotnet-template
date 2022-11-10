FROM mcr.microsoft.com/dotnet/aspnet:7.0 AS base
WORKDIR /app
EXPOSE 80
EXPOSE 443

FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build
WORKDIR /src
COPY ["Template.WebApplication/Template.WebApplication.csproj", "Template.WebApplication/"]
RUN dotnet restore "Template.WebApplication/Template.WebApplication.csproj"
COPY . .
WORKDIR "/src/Template.WebApplication"
RUN dotnet build "Template.WebApplication.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "Template.WebApplication.csproj" -c Release -o /app/publish

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "Template.WebApplication.dll"]
