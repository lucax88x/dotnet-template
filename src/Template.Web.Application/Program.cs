using FluentValidation;
using Template.Web.Application;
using Template.Web.Application.Routing;
using Template.Web.Common.Extensions;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddOptions<Config>()
    .Bind(builder.Configuration.GetSection(Config.Key))
    .ValidateDataAnnotations()
    .ValidateOnStart();

builder
    // .AddTemplateAuth()
    .AddTemplateLogging()
    .RegisterActions()
    .AddSwaggerWhenDevelopment()
    .AddApiExplorer();

builder.Services
    .AddValidatorsFromAssemblyContaining<Program>();

var app = builder.Build();
app.UseHttpsRedirection();

app.UseSwaggerWhenDevelopment()
    // .UseTemplateAuth()
    .UseTemplateRouting()
    .MapCustomerEndpoints()
    .Run();