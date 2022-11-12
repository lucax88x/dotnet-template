using System.ComponentModel.DataAnnotations;

namespace Template.Web.Application;

public record Config {
    public const string Key = "Config";

    [Required]
    public int Host { get; set; }
}