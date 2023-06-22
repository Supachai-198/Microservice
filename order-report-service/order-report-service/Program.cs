using System.Text;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.IdentityModel.Tokens;
using OrderReportService.Services;

var builder = WebApplication.CreateBuilder(args);
builder.Logging.ClearProviders();
builder.Logging.AddConsole();

builder.Services.AddDbContext<APIContext>(options =>
    options.UseSqlServer(builder.Configuration.GetConnectionString("APIContext") ?? throw new InvalidOperationException("Connection string 'APIContext' not found.")));

// Add services to the container.
// Add JWT
builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme).AddJwtBearer(
                config =>
                {
                    //config.SaveToken = true;

                    config.TokenValidationParameters = new TokenValidationParameters()
                    {
                        ValidateIssuerSigningKey = true,
                        IssuerSigningKey = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(builder.Configuration.GetSection("JWT_KEY").Value!)),
                        ValidateIssuer = false,
                        ValidateAudience = false,
                        // set clockskew to zero so tokens expire exactly at token expiration time (instead of 5 minutes later)
                        ClockSkew = TimeSpan.Zero
                    };

                    //check token validation
                    //config.TokenValidationParameters = new TokenValidationParameters()
                    //{
                    //    ValidateLifetime = true,
                    //    ClockSkew = TimeSpan.Zero,
                    //    //ValidateIssuer = true,
                    //    //ValidIssuer = "pea",
                    //    //ValidateAudience = true,
                    //    //ValidAudience = "everyone",
                    //    ValidateIssuerSigningKey = true,
                    //    IssuerSigningKey = new SymmetricSecurityKey(Encoding.UTF8.GetBytes("odn051PvFMtRTBZsqmWkGJl8CHbKceQz"))
                    //};

                }
            );

builder.Services.AddControllers();
// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

//rabbitmq
builder.Services.AddHostedService<ServiceRabbitMQ>();

var app = builder.Build();


// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

// app.UseHttpsRedirection();

app.UseAuthorization();

app.MapControllers();

app.Run();
