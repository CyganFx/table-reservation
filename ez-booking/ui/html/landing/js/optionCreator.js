const createAnOption = (text, value) => {
  const option = document.createElement("option");
  option.text = text;
  option.value = value;
  return option;
};

const createDisabledOption = () => {
  const disabledOption = document.createElement("option");
  disabledOption.text = "Time";
  disabledOption.disabled = true;
  disabledOption.selected = true;
  return disabledOption;
};

const calendar = document.getElementById("clientDate");
calendar.min = new Date()
  .toLocaleDateString("pt-br")
  .split("/")
  .reverse()
  .join("-");

const time_manager = document.getElementById("time_manager");
time_manager.appendChild(createDisabledOption());

const timeChangeHandler = (date) => {
  const time_manager = document.getElementById("time_manager");
  time_manager.innerHTML = "";
  time_manager.appendChild(createDisabledOption());

  let dateTimeHours = new Date().getHours();
  let dateTimeMinutes = new Date().getMinutes();

  if (dateTimeMinutes >= 20 && dateTimeMinutes <= 50) {
    dateTimeMinutes = 30;
  } else {
    dateTimeMinutes = 0;
  }

  if (
    new Date(calendar.value).toLocaleDateString() >=
    new Date().toLocaleDateString()
  ) {
    available_time = Array(15)
      .fill()
      .map((_, idx) => 9 + idx);

    let flag = true;

    if (
      new Date(calendar.value)
        .toLocaleDateString("pt-br")
        .split("/")
        .reverse()
        .join("-") ===
      new Date().toLocaleDateString("pt-br").split("/").reverse().join("-")
    ) {
      if (available_time.includes(15)) {
        available_time = available_time.filter((item) => item >= 15);
      }

      if (dateTimeMinutes === 30 && 15 > 8) {
        available_time.shift();
      }

      if (dateTimeMinutes === 0 && 15 > 8) {
        flag = false;
      }
    }

    available_time.forEach((element) => {
      let number_of_halfs;
      if (flag === false || element === 23) {
        number_of_halfs = 1;
      } else {
        number_of_halfs = 2;
      }
      for (
        let half = 0, minutes = "00";
        half < number_of_halfs;
        half++, minutes = minutes === "00" ? "30" : "00"
      ) {
        if (flag) {
          time_manager.appendChild(createAnOption(element + ":" + minutes));
        } else {
          time_manager.appendChild(createAnOption(element + ":" + "30"));
          flag = true;
        }
      }
    });
  }
};
