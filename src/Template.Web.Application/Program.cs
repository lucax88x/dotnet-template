using FluentValidation;
using Template.Web.Application;
using Template.Web.Application.Routing;
using Template.Web.Common.Configs;
using Template.Web.Common.Extensions;

var builder = WebApplication.CreateBuilder(args);

builder
    .AddAndValidateOptions<SeqConfig>(SeqConfig.Section)
    .AddAndValidateOptions<JaegerConfig>(JaegerConfig.Section)
    // .AddTemplateAuth()
    .AddTemplateLogging()
    .AddTemplateTracing()
    .RegisterActions()
    .AddSwaggerWhenDevelopment();

// TODO: Problem details

builder.Services
    .AddEndpointsApiExplorer()
    .AddValidatorsFromAssemblyContaining<Program>();

var app = builder.Build();
app.UseHttpsRedirection();

app.UseSwaggerWhenDevelopment()
    // .UseTemplateAuth()
    .UseTemplateRouting()
    .MapCustomerEndpoints()
    .Run();