namespace Template.Web.Application;

public class WeatherForecast
{
    public required DateTime Date { get; set; }

    public required int TemperatureC { get; set; }

    public int TemperatureF => 32 + (int)(TemperatureC / 0.5556);

    public required string? Summary { get; set; }
}