using Serilog;
using Serilog.Sinks.FastConsole;
using Template.Common;
using Template.Web.Common.Configs;

namespace Template.Web.Common;

public static class LoggerConfigurationBuilder {
    public static LoggerConfiguration BuildForHost(LoggerConfiguration? loggerConfiguration = null)
    {
        loggerConfiguration ??= new LoggerConfiguration();

        return loggerConfiguration
            .WriteTo.FastConsole();
    }

    public static LoggerConfiguration BuildForApplication(
            LogOptions options,
            LoggerConfiguration? loggerConfiguration = null
        )
    {
        loggerConfiguration ??= new LoggerConfiguration();

        loggerConfiguration
            .Enrich.FromLogContext()
            .Enrich.WithProperty("source", SourceAndVersion.SourceName)
            .WriteTo.FastConsole();

        if (options.SeqConfig is not null)
        {
            loggerConfiguration.WriteTo.Seq(options.SeqConfig.Host);
        }

        return loggerConfiguration;
    }
}

public class LogOptions {
    public SeqConfig? SeqConfig { get; init; }
}