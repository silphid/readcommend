import React  from "react";
import Autocomplete from '@material-ui/lab/Autocomplete';
import TextField from '@material-ui/core/TextField';
import { useState, useEffect } from "react";

type Error = {
  message: string,
}

type Author = {
  id: number,
  firstName: number,
  lastName: number,
}

export default function Authors() {
  const [error, setError] = React.useState<Error | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [authors, setAuthors] = React.useState<Author[]>([]);

  useEffect(() => {
    fetch("/api/v1/authors")
      .then(res => res.json())
      .then(
        (result) => {
          setIsLoaded(true);
          setAuthors(result);
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
      options={authors!}
      getOptionLabel={(option) => `${option.firstName} ${option.lastName}`}
      renderInput={(params) => (
        <TextField
          {...params}
          variant="standard"
          label="Authors"
          placeholder="Type author name(s)"
        />
      )}
    />
    );
  }
}
