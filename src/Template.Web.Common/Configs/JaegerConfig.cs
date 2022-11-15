using System.ComponentModel.DataAnnotations;

namespace Template.Web.Common.Configs;

public record JaegerConfig {
    public const string Section = "Jaeger";

    [Required]
    public required string Host { get; set; }
}