using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using OrderReportService;

namespace OrderReportService.Controllers
{
    [Route("api/v1/[controller]")]
    [ApiController]
    public class OrderReportController : ControllerBase
    {
        private readonly APIContext _context;

        public OrderReportController(APIContext context)
        {
            _context = context;
        }

        // GET: api/OrderReport
        [Authorize]
        [HttpGet]
        public async Task<ActionResult<IEnumerable<OrderReport>>> GetOrderReport()
        {
          if (_context.OrderReport == null)
          {
              return NotFound();
          }
            return await _context.OrderReport.ToListAsync();
        }

        // GET: api/OrderReport/5
        [HttpGet("{id}")]
        public async Task<ActionResult<OrderReport>> GetOrderReport(int id)
        {
          if (_context.OrderReport == null)
          {
              return NotFound();
          }
            var orderReport = await _context.OrderReport.FindAsync(id);

            if (orderReport == null)
            {
                return NotFound();
            }

            return orderReport;
        }

        // PUT: api/OrderReport/5
        // To protect from overposting attacks, see https://go.microsoft.com/fwlink/?linkid=2123754
        [HttpPut("{id}")]
        public async Task<IActionResult> PutOrderReport(int id, OrderReport orderReport)
        {
            if (id != orderReport.Id)
            {
                return BadRequest();
            }

            _context.Entry(orderReport).State = EntityState.Modified;

            try
            {
                await _context.SaveChangesAsync();
            }
            catch (DbUpdateConcurrencyException)
            {
                if (!OrderReportExists(id))
                {
                    return NotFound();
                }
                else
                {
                    throw;
                }
            }

            return NoContent();
        }

        // POST: api/OrderReport
        // To protect from overposting attacks, see https://go.microsoft.com/fwlink/?linkid=2123754
        [HttpPost]
        public async Task<ActionResult<OrderReport>> PostOrderReport(OrderReport orderReport)
        {
          if (_context.OrderReport == null)
          {
              return Problem("Entity set 'APIContext.OrderReport'  is null.");
          }
            _context.OrderReport.Add(orderReport);
            await _context.SaveChangesAsync();

            return CreatedAtAction("GetOrderReport", new { id = orderReport.Id }, orderReport);
        }

        // DELETE: api/OrderReport/5
        [HttpDelete("{id}")]
        public async Task<IActionResult> DeleteOrderReport(int id)
        {
            if (_context.OrderReport == null)
            {
                return NotFound();
            }
            var orderReport = await _context.OrderReport.FindAsync(id);
            if (orderReport == null)
            {
                return NotFound();
            }

            _context.OrderReport.Remove(orderReport);
            await _context.SaveChangesAsync();

            return NoContent();
        }

        private bool OrderReportExists(int id)
        {
            return (_context.OrderReport?.Any(e => e.Id == id)).GetValueOrDefault();
        }
    }
}
