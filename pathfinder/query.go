package pathfinder

import (
	"context"

	pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"
	pfq "github.com/pyrorhythm/libspot/pathfinder/pfrequest"
	pfs "github.com/pyrorhythm/libspot/pathfinder/pfresponse"
)

func (p *Pathfinder) Query(
	ctx context.Context,
	rq pfq.Request,
) (*pfs.Response, error) {
	return p.makeRequest(ctx, rq)
}

func (p *Pathfinder) Top(
	ctx context.Context,
	rq *pfq.SearchTopRequest,
) (*pfs.SearchV2Top, error) {
	resp, err := p.makeRequest(ctx, rq)
	if err != nil {
		return nil, err
	}

	return resp.SearchV2.ToTop(), nil
}

func (p *Pathfinder) Suggestions(
	ctx context.Context,
	rq *pfq.SearchSuggestionsRequest,
) (*pfs.SearchV2Suggestions, error) {
	resp, err := p.makeRequest(ctx, rq)
	if err != nil {
		return nil, err
	}

	return resp.SearchV2.ToSuggestions(), nil
}

func (p *Pathfinder) Tracks(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.SearchV2Tracks, error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchTracks, rq))
	if err != nil {
		return nil, err
	}

	return resp.SearchV2.ToTracks(), nil
}

func (p *Pathfinder) Albums(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.SearchV2Albums, error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchAlbums, rq))
	if err != nil {
		return nil, err
	}

	return resp.SearchV2.ToAlbums(), nil
}

func (p *Pathfinder) Artists(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.SearchV2Artists, error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchArtists, rq))
	if err != nil {
		return nil, err
	}

	return resp.SearchV2.ToArtists(), nil
}

func (p *Pathfinder) Genres(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.SearchV2Genres, error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchGenres, rq))
	if err != nil {
		return nil, err
	}
	return resp.SearchV2.ToGenres(), nil
}

func (p *Pathfinder) Users(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.SearchV2Users, error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchUsers, rq))
	if err != nil {
		return nil, err
	}
	return resp.SearchV2.ToUsers(), nil
}

func (p *Pathfinder) Playlists(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.SearchV2Playlists, error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchPlaylists, rq))
	if err != nil {
		return nil, err
	}
	return resp.SearchV2.ToPlaylists(), nil
}

func (p *Pathfinder) Podcasts(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.SearchV2Podcasts, error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchPodcasts, rq))
	if err != nil {
		return nil, err
	}
	return resp.SearchV2.ToPodcasts(), nil
}

func (p *Pathfinder) Episodes(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.SearchV2Episodes, error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchEpisodes, rq))
	if err != nil {
		return nil, err
	}
	return resp.SearchV2.ToEpisodes(), nil
}

func (p *Pathfinder) Lookup(
	ctx context.Context,
	rq *pfq.SearchSuggestionsRequest,
) ([]*pfd.PlaylistPreviewItems, error) {
	resp, err := p.makeRequest(ctx, rq)
	if err != nil {
		return nil, err
	}

	return resp.Lookup, nil
}

func (p *Pathfinder) GetAlbum(
	ctx context.Context,
	rq *pfq.GetAlbumRequest,
) (*pfd.AlbumFull, error) {
	resp, err := p.makeRequest(ctx, rq)
	if err != nil {
		return nil, err
	}

	return resp.AlbumUnion, nil
}
