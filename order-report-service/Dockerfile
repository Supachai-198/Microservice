FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build-stage
WORKDIR /dotnet-app

COPY *.csproj ./
RUN dotnet restore

COPY ./ ./
RUN dotnet publish -c Release -o dist --no-restore

# production stage
FROM mcr.microsoft.com/dotnet/aspnet:7.0 AS production-stage
WORKDIR /dotnet-app-production
COPY --from=build-stage /dotnet-app/dist ./

EXPOSE 80

ENTRYPOINT [ "dotnet", "OrderReportService.dll" ]