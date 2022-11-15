using System.ComponentModel.DataAnnotations;

namespace Template.Web.Common.Configs;

public record SeqConfig {
    public const string Section = "Seq";

    [Required]
    [Url]
    public required string Host { get; set; }
}