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
        public void Update(T entity);
    }
}
