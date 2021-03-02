import React  from "react";
import Autocomplete from '@material-ui/lab/Autocomplete';
import TextField from '@material-ui/core/TextField';
import { useState, useEffect } from "react";
import * as models from "./models";

export default function Genres(props: { onChange: (genre: models.Genre[]) => void }) {
  const [error, setError] = React.useState<models.Error | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [genres, setGenres] = React.useState<models.Genre[]>([]);

  useEffect(() => {
    fetch("/api/v1/genres")
      .then(res => res.json())
      .then(
        (result) => {
          setIsLoaded(true);
          setGenres(result);
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
      <Autocomplete
      multiple
      options={genres!}
      getOptionLabel={(x) => x.title}
      onChange={(_, items) => props.onChange(items)}
      renderInput={(params) => (
        <TextField
          {...params}
          variant="standard"
          label="Genres"
          placeholder="Type or select genre(s)"
        />
      )}
    />
    );
  }
}
