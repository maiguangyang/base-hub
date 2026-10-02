export const commentContract = {
  method: 'POST',
  pathForCard(cardId) {
    return `/cards/${cardId}/comment-actions`;
  },
  buildBody({ text, plainText = text, atUserIds = [] }) {
    return {
      text,
      plainText,
      atUserIds,
    };
  },
};

export const moveContract = {
  method: 'PATCH',
  pathForCard(cardId) {
    return `/cards/${cardId}`;
  },
  buildBody({ listId, position }) {
    return {
      listId,
      position,
    };
  },
};
