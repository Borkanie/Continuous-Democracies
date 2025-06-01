using ParliamentMonitor.Contracts.Model;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ParliamentMonitor.Contracts.Services
{
    public interface IDBMService<T> where T : Entity
    {
        /// <summary>
        /// Up[dates an entity in the db based on Key.
        /// </summary>
        /// <param name="entity"></param>
        public void Update(T entity);

        /// <summary>
        /// Removes an entity form the DB based on Key.
        /// </summary>
        /// <param name="entity"></param>
        public void Delete(T entity);
    }
}
