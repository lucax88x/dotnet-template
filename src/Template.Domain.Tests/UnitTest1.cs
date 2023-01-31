using FluentAssertions;
using Template.Domain.Customers;

namespace Template.Domain.Tests;

public class UnitTest1 {
    [Fact]
    public void ShouldBeTrue()
    {
        var customer = new Customer() { Name = "name", Id = 1 };

        Assert.Equal("name", customer.Name);
    }
    //
    // [Fact]
    // public void ShouldBeFalse()
    // {
    //     var customer = new Customer() { Name = "name", Id = 1 };
    //     
    //     Assert.Equal("name2", customer.Name);
    // }

    [Fact]
    [Trait("Category", "Integration")]
    public async Task ShouldBeCallSeq()
    {
        using var c = new HttpClient();

        var a = await c.GetAsync(new Uri("http://localhost:5341")).ConfigureAwait(false);

        a.Should().BeSuccessful();
    }
}