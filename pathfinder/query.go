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
) (*pfs.ItemsV2[pfd.OneofMatched], error) {
	resp, err := p.makeRequest(ctx, rq)
	if err != nil {
		return nil, err
	}

	return resp.SearchV2.TopResults, nil
}

func (p *Pathfinder) Suggestions(
	ctx context.Context,
	rq *pfq.SearchSuggestionsRequest,
) (*pfs.ItemsV2[pfd.OneofMatched], error) {
	resp, err := p.makeRequest(ctx, rq)
	if err != nil {
		return nil, err
	}

	return resp.SearchV2.TopResults, nil
}

func (p *Pathfinder) Tracks(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.Items[pfd.Track], error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchTracks, rq))
	if err != nil {
		return nil, err
	}

	return resp.SearchV2.Tracks, nil
}

func (p *Pathfinder) Albums(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.Items[pfd.Album], error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchAlbums, rq))
	if err != nil {
		return nil, err
	}

	return resp.SearchV2.Albums, nil
}

func (p *Pathfinder) Artists(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.Items[pfd.Artist], error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchArtists, rq))
	if err != nil {
		return nil, err
	}

	return resp.SearchV2.Artists, nil
}

func (p *Pathfinder) Genres(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.Items[pfd.Genre], error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchGenres, rq))
	if err != nil {
		return nil, err
	}
	return resp.SearchV2.Genres, nil
}


func (p *Pathfinder) Users(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.Items[pfd.User], error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchUsers, rq))
	if err != nil {
		return nil, err
	}
	return resp.SearchV2.Users, nil
}


func (p *Pathfinder) Playlists(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.Items[pfd.Playlist], error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchPlaylists, rq))
	if err != nil {
		return nil, err
	}
	return resp.SearchV2.Playlists, nil
}


func (p *Pathfinder) Podcasts(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.Items[pfd.Podcast], error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchPodcasts, rq))
	if err != nil {
		return nil, err
	}
	return resp.SearchV2.Podcasts, nil
}

func (p *Pathfinder) Episodes(
	ctx context.Context,
	rq *pfq.BadgeRequestOpts,
) (*pfs.Items[pfd.Episode], error) {
	resp, err := p.makeRequest(ctx, pfq.BadgeSearchFromOpts(pfq.OpSearchEpisodes, rq))
	if err != nil {
		return nil, err
	}
	return resp.SearchV2.Episodes, nil
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
