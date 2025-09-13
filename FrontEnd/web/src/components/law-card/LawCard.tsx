import styles from './LawCard.module.css';

const { Div, title, description, footer, status, date, icon, header } = styles;

export const LawCard = () => {
  return (
    <div className={Div}>
      <div className={header}>
        <h3 className={title}>Climate Change Mitigation Act 2024</h3>
        {/* <div className={status}>
          <div className={icon}></div>
          <p>Activ</p>
        </div> */}
      </div>

      <p className={description}>
        Comprehensive legislation to reduce carbon emissions by 50% by 2030
      </p>
      <div className={footer}>
        <p className={date}>25.11.2024</p>
        <div className={status}>
          <div className={icon}></div>
          <p>Activa</p>
        </div>
      </div>
    </div>
  );
};
