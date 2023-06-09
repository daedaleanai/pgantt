
// Deserialises a value from the local storage.
// Returns null when the value is missing or the deserialisation fails.
export const loadState = (name) => {
  try {
    const serialState = localStorage.getItem(name);
    return JSON.parse(serialState);
  } catch (err) {
    return null;
  }
};

// Serializes a value into the local storage.
export const saveState = (name, state) => {
  try {
    const serialState = JSON.stringify(state);
    localStorage.setItem(name, serialState);
  } catch(err) {
    console.log(err);
  }
};

