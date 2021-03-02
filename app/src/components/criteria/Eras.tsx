import React  from "react";
import { useState, useEffect } from "react";
import * as models from "./models";
import { MenuItem, Select } from "@material-ui/core";

export default function Eras(props: { onChange: (era: models.Era) => void }) {
  const [error, setError] = React.useState<models.Error | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [eras, setEras] = React.useState<models.Era[]>([]);

  useEffect(() => {
    fetch("/api/v1/eras")
      .then(res => res.json())
      .then(
        (result) => {
          setIsLoaded(true);
          setEras(result);
        },
        (error) => {
          setIsLoaded(true);
          setError(error);
        }
      )
  }, [])

  if (error) {
    return <div>Error: {error.message}</div>;
  } else if (!isLoaded) {
    return <div>Loading...</div>;
  } else {
    return (
      <Select title="Title">
        {eras.map((x, i) =>
          <MenuItem key={x.id} onChange={() => { console.log("Selected era:", x.title); props.onChange(x); }}>
            {x.title}
          </MenuItem>
        )}
      </Select>
    );
  }
}
