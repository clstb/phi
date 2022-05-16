import React, {useState} from "react";

export function useInput() {
  const [value, setValue] = useState<string>();
  return {
    value,
    setValue,
    bind: {
      value,
      onChange: (event: React.ChangeEvent<HTMLInputElement>) => {
        setValue(event.target.value);
      }
    }
  };
}

