using FluentValidation;
using Template.WebApplication;
using Template.WebApplication.Routing;

var builder = WebApplication.CreateBuilder(args);

// builder.Services
//     .AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
//     .AddMicrosoftIdentityWebApi(builder.Configuration.GetSection("AzureAd"));

builder.SetLogging()
    .RegisterActions()
    .SetApiExplorerAndSwagger();

builder.Services
    .AddValidatorsFromAssemblyContaining<Program>();

var app = builder.Build();
app.UseHttpsRedirection();

app.SetSwaggerWhenDevelopment()
    // .SetAuth()
    .SetRouting()
    .MapCustomerEndpoints()
    .Run();
