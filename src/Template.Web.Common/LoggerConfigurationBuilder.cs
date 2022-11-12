using Serilog;
using Serilog.Sinks.FastConsole;

namespace Template.Web.Common;

public static class LoggerConfigurationBuilder {
    public static LoggerConfiguration BuildForHost(LoggerConfiguration? loggerConfiguration = null)
    {
        loggerConfiguration ??= new LoggerConfiguration();

        return loggerConfiguration
            .WriteTo.FastConsole();
    }

    public static LoggerConfiguration BuildForApplication(LoggerConfiguration? loggerConfiguration = null)
    {
        loggerConfiguration ??= new LoggerConfiguration();

        loggerConfiguration
            .Enrich.FromLogContext()
            .WriteTo.FastConsole();

        // var configuration = services.GetRequiredService<IConfiguration>();
        // var seqHost = configuration["Seq:Host"] ?? "localhost";
        //
        // loggerConfiguration.WriteTo.Seq($"http://{seqHost}:5341");

        return loggerConfiguration;
    }
}