using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ParliamentMonitor.Contracts.Model
{
    public class Entity
    {
        /// <summary>
        /// Unique ID of the entity, will be used to identyfy the entity.
        /// </summary>
        [Key]
        public required Guid Id { get; set; }

        /// <summary>
        /// The Full Name of the Entity.
        /// </summary>
        public String Name { get; set; } = String.Empty;

        public override bool Equals(object? obj)
        {

            return obj != null && GetType() == obj.GetType() && ((Entity)obj).Id == Id;
        }
    }
}
