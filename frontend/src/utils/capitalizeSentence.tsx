const capitalizeSentence = (s: string) => {
  const t = s
    .split(' ')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
  return t
}

export default capitalizeSentence
