using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;

namespace OrderReportService.Controllers
{
    // localhost:7000/api/v1
    [Route("api/v1")]
    [ApiController]
    public class HomeController : ControllerBase
    {
        [HttpGet]
        public IActionResult Index() {
            return Ok(new { message = "Order Report Service..." });
        }
    }
}
