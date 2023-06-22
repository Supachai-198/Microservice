using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;
using System.Text.Json.Serialization;

namespace OrderReportService {
    public class OrderReport {

        [DatabaseGenerated(DatabaseGeneratedOption.Identity)]
        public int Id { get; set; }

        [Required]
        [JsonPropertyName("user_id")]
        public int UserId { get; set; }

        [JsonPropertyName("product_id")]
        public int ProductId { get; set; }
    }
}